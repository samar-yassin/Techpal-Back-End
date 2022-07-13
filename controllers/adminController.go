package controllers

import (
	"CareerGuidance/database"
	"CareerGuidance/models"
	"context"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	gomail "gopkg.in/gomail.v2"
)

var TrackCollection *mongo.Collection = database.OpenCollection(database.Client, "tracks")

func AddTrack() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var track models.Track
		if err := c.BindJSON(&track); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		validationErr := validate.Struct(track)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": validationErr})
			return
		}

		track.ID = primitive.NewObjectID()
		track.Track_id = track.ID.Hex()

		_, err := TrackCollection.InsertOne(ctx, track)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, track)
	}
}

func AcceptMentor() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var email map[string]string
		var mentor models.Mentor
		if err := c.BindJSON(&email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		err := mentorsCollection.FindOneAndUpdate(ctx, bson.M{"email": email["email"]}, bson.M{"$set": bson.M{"accepted": true}}).Decode(&mentor)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "email or password is incorrect"})
			return
		}

		if mentor.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "mentor not found"})
			return
		}

		subject := "Welcome To Techpal"
		body := "Congratulations! You have been accepted as a mentor in Techpal! <br> You can now login to you account, your password is : " + *mentor.Password

		password, err := bcrypt.GenerateFromPassword([]byte(*mentor.Password), 15)
		if err != nil {
			log.Println(err)
			return
		}
		var tempPass = string(password)

		err = mentorsCollection.FindOneAndUpdate(ctx, bson.M{"email": email["email"]}, bson.M{"$set": bson.M{"password": tempPass}}).Decode(&mentor)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "email or password is incorrect"})
			return
		}

		msg := gomail.NewMessage()
		msg.SetHeader("From", "techpal.guidance@gmail.com")
		msg.SetHeader("To", *mentor.Email)
		msg.SetHeader("Subject", subject)
		msg.SetBody("text/html", body)

		n := gomail.NewDialer("smtp.gmail.com", 587, "techpal.guidance@gmail.com", "osijygroequycuww")

		// Send the email
		if err := n.DialAndSend(msg); err != nil {
			panic(err)
		}

	}
}

func GetAcceptedMentors() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		cursor, err := mentorsCollection.Find(ctx, bson.M{})
		if err != nil {
			log.Println(err)
		}
		defer cursor.Close(ctx)
		var mentors []models.Mentor
		for cursor.Next(ctx) {
			var user models.Mentor
			if err = cursor.Decode(&user); err != nil {
				log.Println(err)
			}
			if *user.User_type == "mentor" && user.Accepted {
				mentors = append(mentors, user)
			}
		}
		c.JSON(http.StatusOK, mentors)
	}
}

func GetNotAcceptedMentors() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		cursor, err := mentorsCollection.Find(ctx, bson.M{})
		if err != nil {
			log.Println(err)
		}
		defer cursor.Close(ctx)
		var mentors []models.Mentor
		for cursor.Next(ctx) {
			var user models.Mentor
			if err = cursor.Decode(&user); err != nil {
				log.Println(err)
			}
			if *user.User_type == "mentor" && !user.Accepted {
				mentors = append(mentors, user)
			}
		}
		c.JSON(http.StatusOK, mentors)
	}
}

func RemoveMentor() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		result, err := mentorsCollection.DeleteOne(ctx, bson.M{"user_id": userId})
		if err != nil {
			log.Println(err)
		}
		var message string
		if result.DeletedCount < 1 {
			message = userId + " doesn't exist."
		} else {
			message = userId + " deleted successfully."
		}
		c.JSON(http.StatusOK, message)
	}
}
