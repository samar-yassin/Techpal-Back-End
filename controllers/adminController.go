package controllers

import (
	"CareerGuidance/database"
	"CareerGuidance/models"
	"context"
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

		err := userCollection.FindOneAndUpdate(ctx, bson.M{"email": email["email"]}, bson.M{"$set": bson.M{"accepted": true}}).Decode(&mentor)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "email or password is incorrect"})
			return
		}

		if mentor.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "mentor not found"})
			return
		}

		subject := "Congratulations"
		body := "Your password is : " + *mentor.Password

		msg := gomail.NewMessage()
		msg.SetHeader("From", "from@gmail.com")
		msg.SetHeader("To", *mentor.Email)
		msg.SetHeader("Subject", subject)
		msg.SetBody("text/html", body)

		n := gomail.NewDialer("smtp.gmail.com", 587, "from@gmail.com", "password")

		// Send the email
		if err := n.DialAndSend(msg); err != nil {
			panic(err)
		}

	}
}