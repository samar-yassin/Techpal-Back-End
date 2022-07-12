package controllers

import (
	"CareerGuidance/database"
	"CareerGuidance/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var SessionsCollection *mongo.Collection = database.OpenCollection(database.Client, "sessions")

func AddSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var session models.Session
		if err := c.BindJSON(&session); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		validationErr := validate.Struct(session)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": validationErr})
			return
		}

		session.ID = primitive.NewObjectID()
		session.SessionId = session.ID.Hex()
		session.MentorId = userId

		_, err := SessionsCollection.InsertOne(ctx, session)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, session)
	}
}

func GetAllSessionsForMentor() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		cursor, err := SessionsCollection.Find(ctx, bson.M{"mentorid": userId})
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

func RemoveSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId := c.Param("session_id")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		result, err := SessionsCollection.DeleteOne(ctx, bson.M{"sessionid": sessionId})
		if err != nil {
			log.Println(err)
		}
		var message string
		if result.DeletedCount < 1 {
			message = sessionId + " doesn't exist."
		} else {
			message = sessionId + " deleted successffuly."
		}
		c.JSON(http.StatusOK, message)
	}
}
