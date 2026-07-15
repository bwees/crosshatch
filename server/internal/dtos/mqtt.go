package dtos

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Some bambu fields arrive as numeric strings, but we want them as numbers
type Number float64

func (n *Number) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == "null" || s == `""` {
		return nil
	}

	// Quoted numeric string: unwrap and parse.
	if len(s) > 0 && s[0] == '"' {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}
		str = strings.TrimSpace(str)
		if str == "" {
			return nil
		}
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil
		}
		*n = Number(f)
		return nil
	}

	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	*n = Number(f)
	return nil
}

type BambuPrintState struct {
	GcodeState             string       `json:"gcode_state"`
	PrintError             Number       `json:"print_error"`
	StgCur                 *Number      `json:"stg_cur"`
	McPercent              Number       `json:"mc_percent"`
	File                   string       `json:"file"`
	McRemainingTime        *Number      `json:"mc_remaining_time"`
	BedTemper              Number       `json:"bed_temper"`
	BedTargetTemper        Number       `json:"bed_target_temper"`
	NozzleTemper           Number       `json:"nozzle_temper"`
	NozzleTargetTemper     Number       `json:"nozzle_target_temper"`
	SpdLvl                 *Number      `json:"spd_lvl"`
	CoolingFanSpeed        Number       `json:"cooling_fan_speed"`
	BigFan1Speed           Number       `json:"big_fan1_speed"`
	BigFan2Speed           Number       `json:"big_fan2_speed"`
	SupportChamberTempEdit bool         `json:"support_chamber_temp_edit"`
	Fun2                   string       `json:"fun2"`
	Device                 BambuDevice  `json:"device"`
	AMS                    BambuAMS     `json:"ams"`
	VirSlot                []BambuTray  `json:"vir_slot"`
	LightsReport           []BambuLight `json:"lights_report"`
}

type BambuDevice struct {
	Ctc struct {
		Info struct {
			// Temp packs the current chamber temperature in its low 16 bits
			// and the target in its high 16 bits.
			Temp Number `json:"temp"`
		} `json:"info"`
	} `json:"ctc"`
}

type BambuAMS struct {
	TrayNow string         `json:"tray_now"`
	AMS     []BambuAMSUnit `json:"ams"`
}

type BambuAMSUnit struct {
	ID          Number      `json:"id"`
	Info        string      `json:"info"`
	HumidityRaw Number      `json:"humidity_raw"`
	DryTime     Number      `json:"dry_time"`
	Temp        Number      `json:"temp"`
	Tray        []BambuTray `json:"tray"`
}

type BambuTray struct {
	ID            *Number `json:"id"`
	TrayType      string  `json:"tray_type"`
	TraySubBrands string  `json:"tray_sub_brands"`
	TrayInfoIdx   string  `json:"tray_info_idx"`
	TrayColor     string  `json:"tray_color"`
	NozzleTempMin *Number `json:"nozzle_temp_min"`
	NozzleTempMax *Number `json:"nozzle_temp_max"`
	Remain        *Number `json:"remain"`
}

type BambuLight struct {
	Node string `json:"node"`
	Mode string `json:"mode"`
}
