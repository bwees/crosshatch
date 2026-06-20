package database

import (
	"crosshatch/internal/database/models"

	"go.uber.org/fx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("crosshatch.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Printer{})

	return db
}

var Module = fx.Provide(
	NewDatabase,
)
