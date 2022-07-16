package controllers

import (
	"CareerGuidance/database"
	"CareerGuidance/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var RatingsCollection *mongo.Collection = database.OpenCollection(database.Client, "ratings")

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		userId := c.Param("user_id")
		var user models.Student
		var mentor models.Mentor
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil || user.Email == nil {
			err = mentorsCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&mentor)
			if err != nil || mentor.Email == nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, mentor)
		} else {
			c.JSON(http.StatusOK, user)
		}
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

func ChangeMentorPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var userpassword map[string]string
		if err := c.BindJSON(&userpassword); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		var mentor models.Mentor
		err := mentorsCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&mentor)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		passwordIsValid, msg := VerifyPassword(*mentor.Password, userpassword["current_password"])
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
		mentor.Password = &tempPass

		err = mentorsCollection.FindOneAndUpdate(ctx, bson.M{"user_id": userId}, bson.M{"$set": bson.M{"password": mentor.Password}}).Decode(&mentor)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, mentor)
	}

}

func GetAllSessions() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		findOptions := options.Find()
		findOptions.SetSort(bson.D{{"date", -1}})
		cursor, err := SessionsCollection.Find(ctx, bson.M{}, findOptions)
		if err != nil {
			log.Println(err)
		}
		defer cursor.Close(ctx)
		var sessions []models.Session
		for cursor.Next(ctx) {
			var session models.Session
			if err = cursor.Decode(&session); err != nil {
				log.Println(err)
			}
			sessions = append(sessions, session)
		}
		c.JSON(http.StatusOK, sessions)
	}
}

func GetEnrolledCourses() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		profileId := c.Param("profile_id")
		defer cancel()
		cursor, err := EnrolledCoursesCollection.Find(ctx, bson.M{"profile_id": profileId, "completed": false})
		if err != nil {
			log.Println(err)
		}
		defer cursor.Close(ctx)
		var courses []models.EnrolledCourse
		for cursor.Next(ctx) {
			var course models.EnrolledCourse
			var c models.Course
			cursor.Decode(&c)
			err := CoursesCollection.FindOne(ctx, bson.M{"course_id": c.Course_id}).Decode(&course)
			course.Completed = c.Completed
			if err != nil {
				log.Println(err)
			}
			courses = append(courses, course)
		}
		c.JSON(http.StatusOK, courses)
	}
}

func GetCompletedCourses() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		profileId := c.Param("profile_id")
		defer cancel()
		cursor, err := EnrolledCoursesCollection.Find(ctx, bson.M{"profile_id": profileId, "completed": true})
		if err != nil {
			log.Println(err)
		}
		defer cursor.Close(ctx)
		var courses []models.EnrolledCourse
		for cursor.Next(ctx) {
			var course models.EnrolledCourse
			var c models.Course
			cursor.Decode(&c)
			err := CoursesCollection.FindOne(ctx, bson.M{"course_id": c.Course_id}).Decode(&course)
			course.Completed = c.Completed
			if err != nil {
				log.Println(err)
			}
			courses = append(courses, course)
		}
		c.JSON(http.StatusOK, courses)
	}
}

func RateCourse() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var rating models.Rating
		if err := c.BindJSON(&rating); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		rating.ID = primitive.NewObjectID()
		rating.Rating_ID = rating.ID.Hex()

		var student models.Student
		err := userCollection.FindOne(ctx, bson.M{"user_id": rating.User_ID}).Decode(&student)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
		*student.Course_rated = *student.Course_rated + 1
		err = userCollection.FindOneAndUpdate(ctx, bson.M{"user_id": rating.User_ID}, bson.M{"$set": student}).Decode(&student)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		_, err = RatingsCollection.InsertOne(ctx, rating)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, rating)
	}
}
