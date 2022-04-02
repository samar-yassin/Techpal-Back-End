package controllers

import (
	"CareerGuidance/database"
	"CareerGuidance/models"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var TracksCollection *mongo.Collection = database.OpenCollection(database.Client, "tracks")

func GetAllTracks() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		cursor, err := TracksCollection.Find(ctx, bson.M{})
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(ctx)
		var tracks []models.Track
		for cursor.Next(ctx) {
			var track models.Track
			if err = cursor.Decode(&track); err != nil {
				log.Fatal(err)
			}
			tracks = append(tracks, track)
		}
		c.JSON(http.StatusOK, tracks)
	}
}
