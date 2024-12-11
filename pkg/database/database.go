package database

import (
	"finance-tracker/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func NewDatabase() *Database {
	conf := config.LoadConfig()

	db, err := gorm.Open(postgres.Open(conf.Database.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed connection to database: %v", err)
	}
	return &Database{db}
}
