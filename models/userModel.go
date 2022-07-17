package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Full_name *string            `json:"full_name" validate:"required,min=2,max=50"`
	User_type *string
	Email     *string
	User_id   string
	Password  *string
	Accepted  bool
}
