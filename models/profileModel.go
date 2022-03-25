package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Profile struct {
	ID               primitive.ObjectID `bson:"_id"`
	Profile_id       string
	Track_id         string   `json:"Track_id"`
	Points           int      `json:"points"`
	Level            int      `json:"level"`
	Completed_skills []string `json:"completed_skills"`
}
