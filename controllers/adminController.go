package controllers

import (
	"CareerGuidance/database"
	"CareerGuidance/models"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
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
