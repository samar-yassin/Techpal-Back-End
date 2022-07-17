package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Report struct {
	ID           primitive.ObjectID `bson:"_id"`
	Report_id    string
	Mentor_email string
	Message      *string `json:"message" validate:"required,min=2,max=300"`
}
