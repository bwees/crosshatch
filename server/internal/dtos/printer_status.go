package dtos

import (
	"math"
	"strconv"
)

// PrinterStage mirrors the numeric stage codes reported by the printer.
type PrinterStage int

type Temperature struct {
	Temperature       float64 `json:"temperature" validate:"required"`
	TargetTemperature float64 `json:"targetTemperature" validate:"required"`
}

type ChamberStatus struct {
	Temperature       float64 `json:"temperature" validate:"required"`
	TargetTemperature float64 `json:"targetTemperature" validate:"required"`
	Controllable      bool    `json:"controllable" validate:"required"`
}

type AMSTray struct {
	ID            int      `json:"id" validate:"required"`
	Empty         bool     `json:"empty" validate:"required"`
	Loaded        bool     `json:"loaded" validate:"required"`
	Material      *string  `json:"material,omitempty"`
	Brand         *string  `json:"brand,omitempty"`
	TrayInfoIdx   *string  `json:"trayInfoIdx,omitempty"`
	Color         *string  `json:"color,omitempty"`
	KValue        *float64 `json:"kValue,omitempty"`
	NozzleTempMin *float64 `json:"nozzleTempMin,omitempty"`
	NozzleTempMax *float64 `json:"nozzleTempMax,omitempty"`
	Remaining     *float64 `json:"remaining,omitempty"`
}

type AMSUnit struct {
	ID             int       `json:"id" validate:"required"`
	Humidity       float64   `json:"humidity" validate:"required"`
	Temperature    float64   `json:"temperature" validate:"required"`
	DryingTime     float64   `json:"dryingTime" validate:"required"`
	SupportsDrying bool      `json:"supportsDrying" validate:"required"`
	Drying         bool      `json:"drying" validate:"required"`
	Trays          []AMSTray `json:"trays" validate:"required"`
}

type FanStatus struct {
	Part    int `json:"part" validate:"required"`
	Aux     int `json:"aux" validate:"required"`
	Chamber int `json:"chamber" validate:"required"`
}

type PrinterStatus struct {
	State         string        `json:"state" validate:"required"`
	Stage         *PrinterStage `json:"stage,omitempty"`
	Progress      float64       `json:"progress" validate:"required"`
	FileName      *string       `json:"fileName,omitempty"`
	TimeRemaining *float64      `json:"timeRemaining,omitempty"`
	BuildPlate    Temperature   `json:"buildPlate" validate:"required"`
	Nozzle        Temperature   `json:"nozzle" validate:"required"`
	Chamber       ChamberStatus `json:"chamber" validate:"required"`
	AMS           []AMSUnit     `json:"ams" validate:"required"`
	ExternalSpool *AMSTray      `json:"externalSpool,omitempty"`
	ChamberLight  bool          `json:"chamberLight" validate:"required"`
	SpeedLevel    *int          `json:"speedLevel,omitempty"`
	Fans          FanStatus     `json:"fans" validate:"required"`
}

func StatusFromMQTT(s *BambuPrintState) PrinterStatus {
	if s == nil {
		return PrinterStatus{AMS: []AMSUnit{}}
	}

	status := PrinterStatus{
		State:        s.GcodeState,
		Progress:     float64(s.McPercent),
		BuildPlate:   Temperature{Temperature: float64(s.BedTemper), TargetTemperature: float64(s.BedTargetTemper)},
		Nozzle:       Temperature{Temperature: float64(s.NozzleTemper), TargetTemperature: float64(s.NozzleTargetTemper)},
		AMS:          []AMSUnit{},
		ChamberLight: chamberLightOn(s.LightsReport),
		Fans: FanStatus{
			Part:    fanPercent(s.CoolingFanSpeed),
			Aux:     fanPercent(s.BigFan1Speed),
			Chamber: fanPercent(s.BigFan2Speed),
		},
	}

	if s.SpdLvl != nil {
		level := int(*s.SpdLvl)
		status.SpeedLevel = &level
	}

	if s.StgCur != nil {
		stage := PrinterStage(int(*s.StgCur))
		status.Stage = &stage
	}

	if s.File != "" {
		file := s.File
		status.FileName = &file
	}

	if s.McRemainingTime != nil {
		remaining := float64(*s.McRemainingTime)
		status.TimeRemaining = &remaining
	}

	// Chamber temperature is packed into a single value: the low 16 bits hold
	// the current temperature, the high 16 bits the target.
	ctcTemp := int(s.Device.Ctc.Info.Temp)
	status.Chamber = ChamberStatus{
		Temperature:       float64(ctcTemp & 0xffff),
		TargetTemperature: float64((ctcTemp >> 16) & 0xffff),
		Controllable:      s.SupportChamberTempEdit,
	}

	// Bit 5 of the fun2 feature mask advertises firmware support for remote drying.
	supportRemoteDry := flagBits(s.Fun2, 5, 1) == 1

	for _, unit := range s.AMS.AMS {
		unitID := int(unit.ID)
		// The AMS info hex packs the hardware model in bits [0,4) and the live
		// drying state in bits [4,8). Only AMS 2 Pro (3) and AMS HT (4) can dry.
		amsType := flagBits(unit.Info, 0, 4)
		dryStatus := flagBits(unit.Info, 4, 4)
		amsUnit := AMSUnit{
			ID:             unitID,
			Humidity:       float64(unit.HumidityRaw),
			Temperature:    float64(unit.Temp),
			DryingTime:     float64(unit.DryTime),
			SupportsDrying: supportRemoteDry && (amsType == 3 || amsType == 4),
			Drying:         isDrying(dryStatus),
			Trays:          []AMSTray{},
		}
		for _, tray := range unit.Tray {
			trayID := 0
			if tray.ID != nil {
				trayID = int(*tray.ID)
			}
			globalID := unitID*4 + trayID
			amsUnit.Trays = append(amsUnit.Trays, trayFromMQTT(tray, trayID, s.AMS.TrayNow == strconv.Itoa(globalID)))
		}
		status.AMS = append(status.AMS, amsUnit)
	}

	if len(s.VirSlot) > 0 {
		slot := s.VirSlot[0]
		id := 254
		if slot.ID != nil {
			id = int(*slot.ID)
		}
		spool := trayFromMQTT(slot, id, s.AMS.TrayNow == "254")
		status.ExternalSpool = &spool
	}

	return status
}

func trayFromMQTT(tray BambuTray, id int, loaded bool) AMSTray {
	result := AMSTray{
		ID:     id,
		Empty:  tray.TrayType == "",
		Loaded: loaded,
	}
	if tray.TrayType != "" {
		v := tray.TrayType
		result.Material = &v
	}
	if tray.TraySubBrands != "" {
		v := tray.TraySubBrands
		result.Brand = &v
	}
	if tray.TrayInfoIdx != "" {
		v := tray.TrayInfoIdx
		result.TrayInfoIdx = &v
	}
	if tray.TrayColor != "" {
		v := tray.TrayColor
		result.Color = &v
	}
	if tray.NozzleTempMin != nil {
		v := float64(*tray.NozzleTempMin)
		result.NozzleTempMin = &v
	}
	if tray.NozzleTempMax != nil {
		v := float64(*tray.NozzleTempMax)
		result.NozzleTempMax = &v
	}
	if tray.Remain != nil {
		v := float64(*tray.Remain)
		result.Remaining = &v
	}
	return result
}

// fanPercent converts a Bambu fan gear reading (0-15) into a 0-100 percentage.
func fanPercent(raw Number) int {
	pct := int(math.Round(float64(raw) / 15.0 * 100.0))
	return min(max(pct, 0), 100)
}

// flagBits parses a hex-encoded bitmask string and extracts `count` bits
// starting at `start` (counted from the least significant bit), matching
// BambuStudio's DevUtil::get_flag_bits. Empty or malformed input yields 0.
func flagBits(hex string, start, count int) int {
	if hex == "" {
		return 0
	}
	v, err := strconv.ParseUint(hex, 16, 64)
	if err != nil {
		return 0
	}
	return int((v >> uint(start)) & ((1 << uint(count)) - 1))
}

// isDrying reports whether an AMS drying status code represents an active
// cycle, matching BambuStudio's DevAms::AmsIsDrying (checking/drying/error).
func isDrying(status int) bool {
	switch status {
	case 1, 2, 5, 6:
		return true
	default:
		return false
	}
}

func chamberLightOn(lights []BambuLight) bool {
	for _, light := range lights {
		if light.Node == "chamber_light" {
			return light.Mode == "on"
		}
	}
	return false
}
