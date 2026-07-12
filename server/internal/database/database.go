package database

import (
	"os"

	"crosshatch/internal/database/models"

	"github.com/glebarez/sqlite"
	"go.uber.org/fx"
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

	db.AutoMigrate(
		&models.Printer{},
		&models.Filament{},
		&models.User{},
		&models.Session{},
	)

	seedFilaments(db)

	return db
}

var Module = fx.Provide(
	NewDatabase,
)
