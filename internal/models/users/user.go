package models

import "time"

type User struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	Password       string    `json:"-"` 
	UserImageUri   string    `json:"user_image_uri"`
	CompanyName    string    `json:"company_name"`
	CompanyImageUri string    `json:"company_image_uri"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
