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

var ProfilesCollection *mongo.Collection = database.OpenCollection(database.Client, "profiles")

func CreateProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var profile models.Profile
		if err := c.BindJSON(&profile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		validationErr := validate.Struct(profile)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": validationErr})
			return
		}

		profile.ID = primitive.NewObjectID()
		profile.Profile_id = profile.ID.Hex()

		_, err := TrackCollection.InsertOne(ctx, profile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, profile)
	}
}
