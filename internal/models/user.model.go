package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name     string
	Email    string `gorm:"unique"`
	Password string
	Balance  float64 `gorm:"default:0"`

	Transactions []Transaction `gorm:"foreignKey:UserID"`
}
