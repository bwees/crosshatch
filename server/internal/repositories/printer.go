package repositories

import (
	"crosshatch/internal/database/models"
	"crosshatch/internal/dtos"

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

func (r *PrinterRepository) GetPrinterBySerial(serial string) (*models.Printer, error) {
	printer := models.Printer{}

	err := r.db.First(&printer, "serial = ?", serial).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &printer, nil
}

func (r *PrinterRepository) CreatePrinter(dto dtos.CreatePrinterDto) (*models.Printer, error) {
	printer := models.Printer{
		Serial:     dto.Serial,
		Name:       dto.Name,
		HostIP:     dto.HostIP,
		AccessCode: dto.AccessCode,
	}

	if err := r.db.Create(&printer).Error; err != nil {
		return nil, err
	}

	return &printer, nil
}

func (r *PrinterRepository) UpdatePrinter(serial string, dto dtos.UpdatePrinterDto) (*models.Printer, error) {
	updates := map[string]interface{}{}
	if dto.Name != nil {
		updates["name"] = *dto.Name
	}
	if dto.HostIP != nil {
		updates["host_ip"] = *dto.HostIP
	}
	if dto.AccessCode != nil {
		updates["access_code"] = *dto.AccessCode
	}

	if len(updates) > 0 {
		err := r.db.Model(&models.Printer{}).Where("serial = ?", serial).Updates(updates).Error
		if err != nil {
			return nil, err
		}
	}

	return r.GetPrinterBySerial(serial)
}

func (r *PrinterRepository) DeletePrinter(serial string) error {
	return r.db.Delete(&models.Printer{}, "serial = ?", serial).Error
}

func NewPrinterRepository(db *gorm.DB) *PrinterRepository {
	return &PrinterRepository{db: db}
}
