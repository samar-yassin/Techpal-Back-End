package controllers

import (
	"CareerGuidance/database"
	"CareerGuidance/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var TracksCollection *mongo.Collection = database.OpenCollection(database.Client, "tracks")
var SkillCollection *mongo.Collection = database.OpenCollection(database.Client, "skills")

func GetAllTracks() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		cursor, err := TracksCollection.Find(ctx, bson.M{})
		if err != nil {
			log.Println(err)
		}
		defer cursor.Close(ctx)
		var tracks []models.Track
		for cursor.Next(ctx) {
			var track models.Track
			if err = cursor.Decode(&track); err != nil {
				log.Println(err)
			}
			tracks = append(tracks, track)
		}
		c.JSON(http.StatusOK, tracks)
	}
}

func GetTrack() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		trackId := c.Param("track_id")
		var track models.Track
		err := TracksCollection.FindOne(ctx, bson.M{"track_id": trackId}).Decode(&track)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, track)

	}
}

func DeleteTrack() gin.HandlerFunc {
	return func(c *gin.Context) {
		var trackId map[string]string
		if err := c.BindJSON(&trackId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		result, err := TracksCollection.DeleteOne(ctx, bson.M{"track_id": trackId["track_id"]})
		if err != nil {
			log.Println(err)
		}
		var message string
		fmt.Println(result.DeletedCount)
		if result.DeletedCount < 1 {
			message = trackId["track_id"] + " doesn't exist."
		} else {
			message = trackId["track_id"] + " deleted successffuly."
		}
		c.JSON(http.StatusOK, message)
	}
}

func AddSkill() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		skill_name := c.Param("skill_name")
		var skill models.Skill

		skill.ID = primitive.NewObjectID()
		skill.Skill_id = skill.ID.Hex()
		skill.Name = skill_name
		_, err := SkillCollection.InsertOne(ctx, skill)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, skill)
	}
}

func GetAllSkills() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		cursor, err := SkillCollection.Find(ctx, bson.M{})
		if err != nil {
			log.Println(err)
		}
		defer cursor.Close(ctx)
		var skills []models.Skill
		for cursor.Next(ctx) {
			var skill models.Skill
			if err = cursor.Decode(&skill); err != nil {
				log.Println(err)
			}
			skills = append(skills, skill)
		}
		c.JSON(http.StatusOK, skills)
	}
}
