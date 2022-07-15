package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Rating struct {
	ID        primitive.ObjectID `bson:"_id"`
	Rating_ID string
	User_ID   string `json:"user_id"`
	Course_ID string `json:"course_id"`
	Rating    int    `json:"rating"`
}
