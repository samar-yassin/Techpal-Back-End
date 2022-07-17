package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Skill struct {
	ID       primitive.ObjectID `bson:"_id"`
	Skill_id string
	Name     string `json:"name"`
}
