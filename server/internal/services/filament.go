package services

import (
	"crosshatch/internal/database/models"
	"crosshatch/internal/dtos"
	"crosshatch/internal/repositories"
)

type FilamentService struct {
	filamentRepo *repositories.FilamentRepository
}

func (s *FilamentService) GetFilaments() ([]models.Filament, error) {
	return s.filamentRepo.GetFilaments()
}

// filamentsFromStatus derives catalog entries from the loaded trays of a
// printer status, so filaments the printer knows about but we don't get added.
func filamentsFromStatus(status dtos.PrinterStatus) []models.Filament {
	filaments := []models.Filament{}

	collect := func(tray *dtos.AMSTray) {
		if tray == nil || tray.TrayInfoIdx == nil || *tray.TrayInfoIdx == "" {
			return
		}

		material := ""
		if tray.Material != nil {
			material = *tray.Material
		}

		name := material
		brand := "Generic"
		if tray.Brand != nil && *tray.Brand != "" {
			name = *tray.Brand
			brand = "Bambu"
		}

		filament := models.Filament{
			TrayInfoIdx: *tray.TrayInfoIdx,
			Brand:       brand,
			Name:        name,
			TrayType:    material,
		}
		if tray.NozzleTempMin != nil {
			filament.NozzleTempMin = *tray.NozzleTempMin
		}
		if tray.NozzleTempMax != nil {
			filament.NozzleTempMax = *tray.NozzleTempMax
		}

		filaments = append(filaments, filament)
	}

	for _, unit := range status.AMS {
		for i := range unit.Trays {
			collect(&unit.Trays[i])
		}
	}
	collect(status.ExternalSpool)

	return filaments
}

func NewFilamentService(filamentRepo *repositories.FilamentRepository) *FilamentService {
	return &FilamentService{filamentRepo: filamentRepo}
}
