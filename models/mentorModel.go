package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Mentor struct {
	ID          primitive.ObjectID `bson:"_id"`
	User_id     string
	Full_name   *string `json:"full_name" validate:"required,min=2,max=50"`
	Email       *string `json:"email" validate:"email,required"`
	Password    *string `json:"password" validate:"required,min=8"`
	User_type   *string `json:"user_type" validate:"required,eq=mentor"`
	Calendly_id *string `json:"calendly_id" validate:"required"`
}
