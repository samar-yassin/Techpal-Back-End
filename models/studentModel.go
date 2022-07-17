package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID              primitive.ObjectID `bson:"_id"`
	User_id         string
	Full_name       *string `json:"full_name" validate:"required,min=2,max=50"`
	Email           *string `json:"email" validate:"email,required"`
	About           *string `json:"about"`
	Degree          *string `json:"degree"`
	Major           *string `json:"major"`
	Start_Year      *string `json:"start_year"`
	End_Year        *string `json:"end_year"`
	Phone           *string `json:"phone"`
	Address         *string `json:"address"`
	Password        *string `json:"password" validate:"required"`
	User_type       *string `json:"user_type" validate:"required,eq=student|eq=mentor"`
	Current_profile string  `json:"current_profile"`
	University      *string `json:"university"`
	Course_rated    *int
	Websites        map[string]string `json:"websites"`
}
