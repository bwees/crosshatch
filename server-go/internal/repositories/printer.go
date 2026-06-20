package repositories

import (
	"crosshatch/internal/database/models"

	"gorm.io/gorm"
)

type PrinterRepository struct {
	db *gorm.DB
}

func (r *PrinterRepository) GetPrinters() ([]models.Printer, error) {
	printers := []models.Printer{}

	err := r.db.Find(&printers).Error
	if err != nil {
		return []models.Printer{}, err
	}

	return printers, nil
}

func NewPrinterRepository(db *gorm.DB) *PrinterRepository {
	return &PrinterRepository{db: db}
}
