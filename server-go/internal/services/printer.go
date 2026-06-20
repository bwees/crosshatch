package services

import (
	"crosshatch/internal/database/models"
	"crosshatch/internal/repositories"
)

type PrinterService struct {
	printerRepo *repositories.PrinterRepository
}

func (s *PrinterService) GetPrinters() ([]models.Printer, error) {
	return s.printerRepo.GetPrinters()
}

func NewPrinterService(printerRepo *repositories.PrinterRepository) *PrinterService {
	return &PrinterService{printerRepo: printerRepo}
}
