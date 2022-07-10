package controllers

import (
	"CareerGuidance/models"
	"context"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		userId := c.Param("user_id")
		var user models.Student
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)

	}
}

func UpdateStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var student models.Student

		userId := c.Param("user_id")
		if err := c.BindJSON(&student); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		err := userCollection.FindOneAndUpdate(ctx, bson.M{"user_id": userId}, bson.M{"$set": bson.M{"full_name": student.Full_name, "email": student.Email, "phone": student.Phone, "address": student.Address, "university": student.University, "websites": student.Websites, "about": student.About, "major": student.Major, "degree": student.Degree, "start_year": student.Start_Year, "end_year": student.End_Year}}).Decode(&student)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, student)

	}
}
func UpdateMentor() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var mentor models.Mentor

		userId := c.Param("user_id")
		if err := c.BindJSON(&mentor); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		err := mentorsCollection.FindOneAndUpdate(ctx, bson.M{"user_id": userId}, bson.M{"$set": bson.M{"full_name": mentor.Full_name, "email": mentor.Email, "calendly_id": mentor.Calendly_id, "about": mentor.About}}).Decode(&mentor)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, mentor)

	}
}

func ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var userpassword map[string]string
		if err := c.BindJSON(&userpassword); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, userpassword["current_password"])
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": msg})
			return
		}

		password, err := bcrypt.GenerateFromPassword([]byte(userpassword["new_password"]), 15)
		if err != nil {
			log.Println(err)
			return
		}
		var tempPass = string(password)
		user.Password = &tempPass

		err = userCollection.FindOneAndUpdate(ctx, bson.M{"user_id": userId}, bson.M{"$set": bson.M{"password": user.Password}}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)

	}

}
