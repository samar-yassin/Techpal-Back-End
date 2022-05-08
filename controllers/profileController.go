package controllers

import (
	"CareerGuidance/database"
	"CareerGuidance/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

		var user models.Student
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
		err := ProfilesCollection.FindOne(ctx, bson.M{"profile_id": profileId["profile_id"]}).Decode(&profile)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		var user models.Student
		err = userCollection.FindOneAndUpdate(ctx, bson.M{"user_id": userId}, bson.M{"$set": bson.M{"current_profile": profileId["profile_id"]}}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, profile)
	}
}

func DeleteProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		var profileId map[string]string
		if err := c.BindJSON(&profileId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		result, err := ProfilesCollection.DeleteOne(ctx, bson.M{"profile_id": profileId["profile_id"]})
		if err != nil {
			log.Fatal(err)
		}
		var message string
		fmt.Println(result.DeletedCount)
		if result.DeletedCount < 1 {
			message = profileId["profile_id"] + " doesn't exist."
		} else {
			message = profileId["profile_id"] + " deleted successffuly."
		}
		c.JSON(http.StatusOK, message)
	}
}

func GetCurrentProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.Student
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, user.Current_profile)
	}
}

func GetAllProfiles() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		cursor, err := ProfilesCollection.Find(ctx, bson.M{})
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(ctx)
		var profiles []models.Profile
		for cursor.Next(ctx) {
			var profile models.Profile
			if err = cursor.Decode(&profile); err != nil {
				log.Fatal(err)
			}
			profiles = append(profiles, profile)
		}
		c.JSON(http.StatusOK, profiles)
	}
}
