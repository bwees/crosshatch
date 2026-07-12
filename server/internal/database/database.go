package database

import (
	"os"
	"path/filepath"

	"crosshatch/internal/config"
	"crosshatch/internal/database/models"

	"github.com/glebarez/sqlite"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func NewDatabase() *gorm.DB {
	dataDir := config.DataDir()
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		panic("failed to create data directory")
	}

	dsn := filepath.Join(dataDir, "crosshatch.db")

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(
		&models.Printer{},
		&models.Filament{},
		&models.User{},
		&models.Session{},
		&models.PushSubscription{},
		&models.NotificationSetting{},
		&models.AppConfig{},
	)

	seedFilaments(db)

	return db
}

var Module = fx.Provide(
	NewDatabase,
)
