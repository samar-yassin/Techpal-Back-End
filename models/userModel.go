package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID              primitive.ObjectID `bson:"_id"`
	User_id         string
	Full_name       *string `json:"full_name" validate:"required,min=2,max=50"`
	Email           *string `json:"email" validate:"email,required"`
	Phone           *string `json:"phone"`
	Address         *string `json:"address"`
	Password        *string `json:"password" validate:"required,min=8"`
	User_type       *string `json:"user_type" validate:"required,eq=Student|Mentor"`
	Current_profile *string `json:"current_profile"`
	University      *string `json:"university"`
}
