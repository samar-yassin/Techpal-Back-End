package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Profile struct {
	ID               primitive.ObjectID `bson:"_id"`
	Profile_id       string
	User_id          string
	Track_id         string   `json:"Track_id" validate:"required"`
	Points           int      `json:"points"`
	Level            int      `json:"level"`
	Completed_Skills []string `json:"completed_skills"`
	EnrolledCourses  []EnrolledCourse
}
