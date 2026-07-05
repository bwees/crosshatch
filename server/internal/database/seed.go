package database

import (
	"crosshatch/internal/database/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// filamentPresets mirrors OrcaSlicer's Bambu profiles. TrayInfoIdx is the
// filament id the printer expects; TrayType is the material family reported to
// the AMS; temperatures are the profile defaults.
var filamentPresets = []models.Filament{
	// Generic
	{Brand: "Generic", Name: "PLA", TrayType: "PLA", TrayInfoIdx: "GFL99", NozzleTempMin: 190, NozzleTempMax: 240},
	{Brand: "Generic", Name: "PLA Silk", TrayType: "PLA", TrayInfoIdx: "GFL96", NozzleTempMin: 190, NozzleTempMax: 240},
	{Brand: "Generic", Name: "PLA-CF", TrayType: "PLA-CF", TrayInfoIdx: "GFL98", NozzleTempMin: 190, NozzleTempMax: 240},
	{Brand: "Generic", Name: "PLA High Speed", TrayType: "PLA", TrayInfoIdx: "GFL95", NozzleTempMin: 190, NozzleTempMax: 240},
	{Brand: "Generic", Name: "PETG", TrayType: "PETG", TrayInfoIdx: "GFG99", NozzleTempMin: 220, NozzleTempMax: 270},
	{Brand: "Generic", Name: "PETG-CF", TrayType: "PETG-CF", TrayInfoIdx: "GFG98", NozzleTempMin: 240, NozzleTempMax: 270},
	{Brand: "Generic", Name: "ABS", TrayType: "ABS", TrayInfoIdx: "GFB99", NozzleTempMin: 240, NozzleTempMax: 280},
	{Brand: "Generic", Name: "ASA", TrayType: "ASA", TrayInfoIdx: "GFB98", NozzleTempMin: 240, NozzleTempMax: 280},
	{Brand: "Generic", Name: "TPU", TrayType: "TPU", TrayInfoIdx: "GFU99", NozzleTempMin: 200, NozzleTempMax: 250},
	{Brand: "Generic", Name: "TPU for AMS", TrayType: "TPU", TrayInfoIdx: "GFU98", NozzleTempMin: 200, NozzleTempMax: 250},
	{Brand: "Generic", Name: "PC", TrayType: "PC", TrayInfoIdx: "GFC99", NozzleTempMin: 260, NozzleTempMax: 290},
	{Brand: "Generic", Name: "PA", TrayType: "PA", TrayInfoIdx: "GFN99", NozzleTempMin: 240, NozzleTempMax: 280},
	{Brand: "Generic", Name: "PA-CF", TrayType: "PA-CF", TrayInfoIdx: "GFN98", NozzleTempMin: 260, NozzleTempMax: 300},
	{Brand: "Generic", Name: "PVA", TrayType: "PVA", TrayInfoIdx: "GFS99", NozzleTempMin: 190, NozzleTempMax: 240},
	{Brand: "Generic", Name: "HIPS", TrayType: "HIPS", TrayInfoIdx: "GFS98", NozzleTempMin: 220, NozzleTempMax: 270},
	{Brand: "Generic", Name: "PP", TrayType: "PP", TrayInfoIdx: "GFP97", NozzleTempMin: 220, NozzleTempMax: 250},

	// Bambu
	{Brand: "Bambu", Name: "PLA Basic", TrayType: "PLA", TrayInfoIdx: "GFA00", NozzleTempMin: 190, NozzleTempMax: 240},
	{Brand: "Bambu", Name: "PLA Matte", TrayType: "PLA", TrayInfoIdx: "GFA01", NozzleTempMin: 190, NozzleTempMax: 240},
	{Brand: "Bambu", Name: "PLA Silk", TrayType: "PLA", TrayInfoIdx: "GFA05", NozzleTempMin: 190, NozzleTempMax: 240},
	{Brand: "Bambu", Name: "PLA Galaxy", TrayType: "PLA", TrayInfoIdx: "GFA15", NozzleTempMin: 190, NozzleTempMax: 240},
	{Brand: "Bambu", Name: "PLA Metal", TrayType: "PLA", TrayInfoIdx: "GFA02", NozzleTempMin: 190, NozzleTempMax: 240},
	{Brand: "Bambu", Name: "PLA-CF", TrayType: "PLA-CF", TrayInfoIdx: "GFA50", NozzleTempMin: 210, NozzleTempMax: 250},
	{Brand: "Bambu", Name: "PETG HF", TrayType: "PETG", TrayInfoIdx: "GFG02", NozzleTempMin: 230, NozzleTempMax: 270},
	{Brand: "Bambu", Name: "PETG Basic", TrayType: "PETG", TrayInfoIdx: "GFG00", NozzleTempMin: 230, NozzleTempMax: 270},
	{Brand: "Bambu", Name: "PETG Translucent", TrayType: "PETG", TrayInfoIdx: "GFG01", NozzleTempMin: 230, NozzleTempMax: 270},
	{Brand: "Bambu", Name: "PETG-CF", TrayType: "PETG-CF", TrayInfoIdx: "GFG50", NozzleTempMin: 240, NozzleTempMax: 270},
	{Brand: "Bambu", Name: "ABS", TrayType: "ABS", TrayInfoIdx: "GFB00", NozzleTempMin: 240, NozzleTempMax: 280},
	{Brand: "Bambu", Name: "ABS-GF", TrayType: "ABS-GF", TrayInfoIdx: "GFB50", NozzleTempMin: 240, NozzleTempMax: 280},
	{Brand: "Bambu", Name: "ASA", TrayType: "ASA", TrayInfoIdx: "GFB01", NozzleTempMin: 240, NozzleTempMax: 280},
	{Brand: "Bambu", Name: "ASA-CF", TrayType: "ASA-CF", TrayInfoIdx: "GFB51", NozzleTempMin: 250, NozzleTempMax: 280},
	{Brand: "Bambu", Name: "TPU 95A", TrayType: "TPU", TrayInfoIdx: "GFU01", NozzleTempMin: 200, NozzleTempMax: 250},
	{Brand: "Bambu", Name: "TPU 95A HF", TrayType: "TPU", TrayInfoIdx: "GFU00", NozzleTempMin: 200, NozzleTempMax: 250},
	{Brand: "Bambu", Name: "PC", TrayType: "PC", TrayInfoIdx: "GFC00", NozzleTempMin: 260, NozzleTempMax: 290},
	{Brand: "Bambu", Name: "PC FR", TrayType: "PC", TrayInfoIdx: "GFC01", NozzleTempMin: 260, NozzleTempMax: 290},
	{Brand: "Bambu", Name: "PA-CF", TrayType: "PA-CF", TrayInfoIdx: "GFN03", NozzleTempMin: 260, NozzleTempMax: 300},
	{Brand: "Bambu", Name: "PAHT-CF", TrayType: "PA-CF", TrayInfoIdx: "GFN04", NozzleTempMin: 260, NozzleTempMax: 300},
	{Brand: "Bambu", Name: "PA6-CF", TrayType: "PA6-CF", TrayInfoIdx: "GFN05", NozzleTempMin: 260, NozzleTempMax: 300},
	{Brand: "Bambu", Name: "PVA", TrayType: "PVA", TrayInfoIdx: "GFS04", NozzleTempMin: 210, NozzleTempMax: 250},
	{Brand: "Bambu", Name: "Support For PLA", TrayType: "PLA", TrayInfoIdx: "GFS02", NozzleTempMin: 190, NozzleTempMax: 240},
	{Brand: "Bambu", Name: "Support For PLA-PETG", TrayType: "PLA", TrayInfoIdx: "GFS05", NozzleTempMin: 190, NozzleTempMax: 240},
	{Brand: "Bambu", Name: "Support W", TrayType: "PLA", TrayInfoIdx: "GFS00", NozzleTempMin: 190, NozzleTempMax: 240},
}

// seedFilaments inserts the built-in filament presets, leaving any existing
// rows (auto-added or edited) untouched.
func seedFilaments(db *gorm.DB) {
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&filamentPresets)
}
