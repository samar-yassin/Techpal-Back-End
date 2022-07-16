package models
import("time")
import "go.mongodb.org/mongo-driver/bson/primitive"

type Session struct {
	ID           primitive.ObjectID `bson:"_id"`
	SessionId    string
	MentorId     string
	Session_Name *string `json:"session_name" validate:"required"`
	Date         *time.Time `json:"date" validate:"required"`
	Meeting_Link *string `json:"meeting_link" validate:"required"`
}
