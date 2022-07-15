package models

type Course struct {
	Profile_id string `json:"profile_id"`
	User_id    string `json:"user_id"`
	Course_id  string `json:"course_id"`
	Completed  *bool  `json:"completed"`
}
