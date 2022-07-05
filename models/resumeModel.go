package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Resume struct {
	ID         primitive.ObjectID `bson:"_id"`
	Profile_id string
	Resume_id  string
	Template   *string      `json:"template" validate:"required"`
	LeftOrder  []resumeUnit `json:"leftorder" validate:"required"`
	RightOrder []resumeUnit `json:"rightorder" validate:"required"`
}
