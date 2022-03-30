package controllers

import (
	"CareerGuidance/database"
	"CareerGuidance/models"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

var ProfilesCollection *mongo.Collection = database.OpenCollection(database.Client, "profiles")

func CreateProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
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
		profile.User_id = userId

		var user models.User
		err := userCollection.FindOneAndUpdate(ctx, bson.M{"user_id": userId}, bson.M{"$set": bson.M{"current_profile": profile.Profile_id}}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user.Current_profile = profile.Profile_id

		_, err = ProfilesCollection.InsertOne(ctx, profile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, profile)
	}
}

func SwitchProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		var profileId map[string]string
		if err := c.BindJSON(&profileId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var profile models.Profile
		err := userCollection.FindOne(ctx, bson.M{"profile_id": profileId}).Decode(&profile)
		log.Println(userId)
		log.Println(profileId["profile_id"])
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		err = userCollection.FindOneAndUpdate(ctx, bson.M{"user_id": userId}, bson.M{"$set": bson.M{"current_profile": profileId["profile_id"]}}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user.Current_profile = profile.Profile_id

		c.JSON(http.StatusOK, profile)
	}
}
