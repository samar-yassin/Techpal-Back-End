package controllers

import (
	"CareerGuidance/database"
	"CareerGuidance/helpers"
	"CareerGuidance/models"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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

func validMailAddress(address string) (string, bool) {
	addr, err := mail.ParseAddress(address)
	if err != nil {
		return "", false
	}
	return addr.Address, true
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

		if &foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "User not found"})
			return
		}

		passwordIsValid, msg := VerifyPassword(user.Password, foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": msg})
			return
		}

		token, refreshToken, _ := helpers.GenerateTokens(&foundUser.Email, &foundUser.First_Name, strconv.FormatUint(uint64(foundUser.User_id), 10))
		helpers.UpdateTokens(token, refreshToken, strconv.FormatUint(uint64(foundUser.User_id), 10))
		c.SetCookie("jwt", token, 60*60*24, "/", "career guidance", true, true)

		c.JSON(http.StatusOK, token)
	}
}

func UploadCv() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		filename := header.Filename
		out, err := os.Create("CVs/" + filename)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"msg": "uploaded successfully"})
	}
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data map[string]string
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		addr_, ok := validMailAddress(data["email"])
		if !(ok) {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Please provide a valide email address"})
			return
		}
		password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 15)
		user := models.User{
			First_Name: data["first_name"],
			Last_Name:  data["last_name"],
			Email:      addr_,
			Password:   string(password),
		}

		// validate := validator.New()
		// err := validate.Struct(user)
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		// 	return
		// }

		podcastResult, err := userCollection.InsertOne(context.TODO(), user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		fmt.Printf("Inserted %v document into users collection!\n", (podcastResult.InsertedID))
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

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("jwt", "", -1, "", "", false, true)
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}
