package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Session struct {
	ID          primitive.ObjectID `bson:"_id"`
	SessionId   string
	MentorId    string
	SessionName *string `json:"session_name" validate:"required"`
	Date        *string `json:"date" validate:"required"`
	Time        *string `json:"time" validate:"required"`
	MeetingLink *string `json:"meeting_link" validate:"required"`
}
