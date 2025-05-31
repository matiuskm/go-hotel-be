package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username 		string		`gorm:"uniqueIndex;not null" json:"username"`
	Password 		string		`gorm:"not null" json:"password"`
	FullName 		string		`json:"full_name"`
	Role 			string		`gorm:"not null" json:"role"`
}