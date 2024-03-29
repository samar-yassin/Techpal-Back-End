package controllers

import (
	"CareerGuidance/database"
	"CareerGuidance/helpers"
	"CareerGuidance/models"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var mentorsCollection *mongo.Collection = database.OpenCollection(database.Client, "mentors")

var validate = validator.New()

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	check := true
	msg := ""
	log.Println(err)
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

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil || foundUser.Email == nil {
			err = mentorsCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
			if err != nil || foundUser.Email == nil || foundUser.Accepted == false {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "not found or not accepted yet"})
			}
		}

		passwordIsValid, msg := VerifyPassword(*foundUser.Password, *user.Password)

		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": msg})
			return
		}

		token, refreshToken, _ := helpers.GenerateTokens(foundUser.Email, foundUser.Full_name, foundUser.User_id)
		helpers.UpdateTokens(token, refreshToken, foundUser.User_id)

		c.SetCookie("jwt", token, 60*60*24, "/", "career guidance", true, true)

		c.JSON(http.StatusOK, token)
		c.JSON(http.StatusOK, foundUser)

	}
}

func ApplyMentor() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var mentor models.Mentor
		if err := c.BindJSON(&mentor); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		rand.Seed(time.Now().UnixNano())
		chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
			"abcdefghijklmnopqrstuvwxyz" +
			"0123456789" +
			"#*$%&!")
		length := 8
		var b strings.Builder
		for i := 0; i < length; i++ {
			b.WriteRune(chars[rand.Intn(len(chars))])
		}
		password := b.String()
		mentor.Password = &password
		var user_type = "mentor"
		mentor.User_type = &user_type

		validationErr := validate.Struct(mentor)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}
		count, err := mentorsCollection.CountDocuments(ctx, bson.M{"email": mentor.Email})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error Occurred"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error Occurred"})
			return
		}

		mentor.ID = primitive.NewObjectID()
		mentor.User_id = mentor.ID.Hex()

		_, err = mentorsCollection.InsertOne(ctx, mentor)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error Occurred while adding user"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, mentor.User_id)

	}
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.Student
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error Occurred"})
			return
		}

		password, err := bcrypt.GenerateFromPassword([]byte(*user.Password), 15)
		if err != nil {
			log.Println(err)
			return
		}
		var tempPass = string(password)
		user.Password = &tempPass

		user.Course_rated = 0

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error Occurred"})
			return
		}

		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, _, _ := helpers.GenerateTokens(user.Email, user.Full_name, user.User_id)

		c.SetCookie("jwt", token, 60*60*24, "/", "career guidance", true, true)

		_, err = userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error Occurred while adding user"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, token)
		c.JSON(http.StatusOK, user)

	}
}

func User() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("jwt")

		if err != nil {
			c.JSON(http.StatusOK, err.Error())
			return
		}
		c.JSON(http.StatusOK, cookie)
	}
}
