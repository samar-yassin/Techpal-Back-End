package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Track struct {
	ID       primitive.ObjectID `bson:"_id"`
	Track_id string
	Name     *string        `json:"name" validate:"required,min=2,max=100"`
	Color1   *string        `json:"color1" validate:"required"`
	Color2   *string        `json:"color2" validate:"required"`
	Skills   map[string]int `json:"skills" validate:"required"`
}
