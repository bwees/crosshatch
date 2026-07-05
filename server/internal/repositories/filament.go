package repositories

import (
	"crosshatch/internal/database/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FilamentRepository struct {
	db *gorm.DB
}

func (r *FilamentRepository) GetFilaments() ([]models.Filament, error) {
	filaments := []models.Filament{}

	err := r.db.Order("brand, name").Find(&filaments).Error
	if err != nil {
		return []models.Filament{}, err
	}

	return filaments, nil
}

// CreateMissing inserts filaments that aren't already in the catalog, leaving
// existing rows untouched.
func (r *FilamentRepository) CreateMissing(filaments []models.Filament) error {
	if len(filaments) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&filaments).Error
}

func NewFilamentRepository(db *gorm.DB) *FilamentRepository {
	return &FilamentRepository{db: db}
}
