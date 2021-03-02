package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name            string `json:"name"`
	Email           string `json:"email" gorm:"uniqueIndex"`
	Password        string `json:"password"`
	IsEmailVerified bool   `json:"is_email_verified"`
}
