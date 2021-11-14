package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Address  string
	Username string
	Phone    string
}
