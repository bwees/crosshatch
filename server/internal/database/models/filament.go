package models

type Filament struct {
	TrayInfoIdx   string  `gorm:"primaryKey" json:"trayInfoIdx" validate:"required"`
	Brand         string  `gorm:"not null" json:"brand" validate:"required"`
	Name          string  `gorm:"not null" json:"name" validate:"required"`
	TrayType      string  `gorm:"not null" json:"trayType" validate:"required"`
	NozzleTempMin float64 `gorm:"not null" json:"nozzleTempMin" validate:"required"`
	NozzleTempMax float64 `gorm:"not null" json:"nozzleTempMax" validate:"required"`
}

func (Filament) TableName() string {
	return "filament"
}
