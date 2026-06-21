package database

import (
	"os"

	"crosshatch/internal/database/models"

	"go.uber.org/fx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "crosshatch.db"
	}

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// err = db.Migrator().RenameTable("printers", &models.Printer{})
	// if err != nil {
	// 	panic("failed to rename table: " + err.Error())
	// }

	db.AutoMigrate(&models.Printer{})

	return db
}

var Module = fx.Provide(
	NewDatabase,
)
