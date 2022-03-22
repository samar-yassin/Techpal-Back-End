package models

type User struct {
	First_Name string `json:"first_name" validate:"required,min=2,max=50"`
	Last_Name  string `json:"last_name" validate:"required,min=2,max=50"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"-" validate:"required,min=8"`
	User_id    uint   `json:"user_id" validate:"omitempty,uuid"`
}
