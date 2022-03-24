package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Full_name *string            `json:"full_name" validate:"required,min=2,max=50"`
	Email     *string            `json:"email" validate:"email,required"`
	Password  *string            `json:"password" validate:"required,min=8"`
	User_id   string
}
