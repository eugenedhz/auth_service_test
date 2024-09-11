package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return &gorm.DB{}, err
	}

	return db, nil
}
