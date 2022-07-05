package controllers

import (
	"CareerGuidance/database"
	"CareerGuidance/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ResumesCollection *mongo.Collection = database.OpenCollection(database.Client, "resumes")

func AddResume() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var resume models.Resume
		if err := c.BindJSON(&resume); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		validationErr := validate.Struct(resume)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": validationErr})
			return
		}

		resume.ID = primitive.NewObjectID()
		resume.Resume_id = resume.ID.Hex()

		_, err := ResumesCollection.InsertOne(ctx, resume)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, resume)
	}
}

func GetResume() gin.HandlerFunc {
	return func(c *gin.Context) {
		profileId := c.Param("profile_id")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var resume models.Resume
		err := ResumesCollection.FindOne(ctx, bson.M{"profile_id": profileId}).Decode(&resume)
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, resume)
	}
}

func UpdateResume() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var resume models.Resume

		profileId := c.Param("profile_id")
		if err := c.BindJSON(&resume); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		err := ResumesCollection.FindOneAndUpdate(ctx, bson.M{"profile_id": profileId}, bson.M{"$set": bson.M{"template": resume.Template, "leftorder": resume.LeftOrder, "rightorder": resume.RightOrder}}).Decode(&resume)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resume)

	}
}
