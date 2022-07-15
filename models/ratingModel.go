package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Rating struct {
	ID        primitive.ObjectID `bson:"_id"`
	RatingID  int
	StudentID string `json:"student_id"`
	CourseID  string `json:"course_id"`
	Rating    int    `json:"rating"`
}
