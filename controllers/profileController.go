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

var CoursesCollection *mongo.Collection = database.OpenCollection(database.Client, "courses")
var EnrolledCoursesCollection *mongo.Collection = database.OpenCollection(database.Client, "enrolledCourses")

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

		var track models.Track
		err := TracksCollection.FindOne(ctx, bson.M{"track_id": profile.Track_id}).Decode(&track)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "No track found"})
			return
		}

		var points int
		points = 0
		//var lvl int
		//lvl = 0
		for _, skill := range profile.Completed_Skills {
			points += track.Skills[skill]
			//lvl++
		}

		profile.Points = points
		//profile.Level = lvl

		var user models.Student
		err = userCollection.FindOneAndUpdate(ctx, bson.M{"user_id": userId}, bson.M{"$set": bson.M{"current_profile": profile.Profile_id}}).Decode(&user)
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
			log.Println(err)
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
			log.Println(err)
		}

		var profile models.Profile
		err2 := ProfilesCollection.FindOne(ctx, bson.M{"profile_id": user.Current_profile}).Decode(&profile)
		if err2 != nil {
			log.Println(err2)
		}
		c.JSON(http.StatusOK, profile)
	}
}

func GetAllProfiles() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		cursor, err := ProfilesCollection.Find(ctx, bson.M{})
		if err != nil {
			log.Println(err)
		}
		defer cursor.Close(ctx)
		var profiles []models.Profile
		for cursor.Next(ctx) {
			var profile models.Profile
			if err = cursor.Decode(&profile); err != nil {
				log.Println(err)
			}
			if profile.User_id == userId {
				profiles = append(profiles, profile)
			}
		}
		c.JSON(http.StatusOK, profiles)
	}
}

func MarkCompleted() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var foundProfile models.Profile
		var profile map[string]string
		var course models.Course

		if err := c.BindJSON(&profile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error1": err.Error()})
			return
		}

		err := ProfilesCollection.FindOne(ctx, bson.M{"profile_id": profile["profile_id"]}).Decode(&foundProfile)

		err = EnrolledCoursesCollection.FindOne(ctx, bson.M{"profile_id": profile["profile_id"], "course_id": profile["course_id"]}).Decode(&course)
		*course.Completed = true

		var track models.Track
		err = TracksCollection.FindOne(ctx, bson.M{"track_id": foundProfile.Track_id}).Decode(&track)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "No track found"})
			return
		}

		var enrolledcourse models.EnrolledCourse
		err = CoursesCollection.FindOne(ctx, bson.M{"course_id": profile["course_id"]}).Decode(&enrolledcourse)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "No course found"})
			return
		}

		foundProfile.Points += track.Skills[*enrolledcourse.Skill]

		//profile.Points = profile.Points+1;
		err = EnrolledCoursesCollection.FindOneAndUpdate(ctx, bson.M{"profile_id": profile["profile_id"], "course_id": profile["course_id"]}, bson.M{"$set": course}).Decode(&course)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error1": err.Error()})
			return
		}

		err = ProfilesCollection.FindOneAndUpdate(ctx, bson.M{"profile_id": profile["profile_id"]}, bson.M{"$set": foundProfile}).Decode(&foundProfile)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error2": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundProfile)

	}
}

func EnrollCourse() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
		var foundProfile models.Profile
		var course models.Course
		if err := c.BindJSON(&course); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		_, err := EnrolledCoursesCollection.InsertOne(ctx, course)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "No course found"})
			return
		}
		/*
			//foundProfile.EnrolledCourses = append(foundProfile.EnrolledCourses, enrolledcourse)

			err = ProfilesCollection.FindOneAndUpdate(ctx, bson.M{"profile_id": profile["profile_id"]}, bson.M{"$enrolledcourses": bson.M{"completed": "false"}}).Decode(&foundProfile)
			defer cancel()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
				return
			}

		*/

		c.JSON(http.StatusOK, foundProfile)

	}
}
