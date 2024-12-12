package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name     string
	Email    string `gorm:"unique"`
	Password string

	Transactions []Transaction `gorm:"foreignKey:UserID"`
}
