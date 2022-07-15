package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type EnrolledCourse struct {
	ID          primitive.ObjectID `bson:"_id"`
	Completed   *bool
	Skill       *string `json:"skill"`
	Course_id   *string `json:"course_id"`
	Course_name *string `json:"course_name"`
	Course_url  *string `json:"course_url"`
}
