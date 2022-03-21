package controllers

import (
	"CareerGuidance/database"
	"CareerGuidance/helpers"
	"CareerGuidance/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("Email or Password is incorrect")
		check = false
	}

	return check, msg
}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		log.Println("email ", &user.Email)

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "email or password is incorrect"})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "User not found"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": msg})
			return
		}

		token, refreshToken, _ := helpers.GenerateTokens(foundUser.Email, foundUser.First_name, foundUser.User_id)
		helpers.UpdateTokens(token, refreshToken, foundUser.User_id)

		c.SetCookie("jwt", token, 60*60*24, "/", "career guidance", true, true)

		c.JSON(http.StatusOK, token)
	}
}
