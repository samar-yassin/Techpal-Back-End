package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ContactUs struct {
	ID         primitive.ObjectID `bson:"_id"`
	First_name string             `json:"first_name"`
	Last_name  string             `json:"last_name"`
	Email      string             `json:"email"`
	Message    string             `json:"message"`
}
