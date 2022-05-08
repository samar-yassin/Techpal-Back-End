package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Report struct {
	ID         primitive.ObjectID `bson:"_id"`
	Report_id  string
	Message    *string `json:"message" validate:"required,min=2,max=300"`
	Session_id *string `json:"session_id" validate:"required"`
}
