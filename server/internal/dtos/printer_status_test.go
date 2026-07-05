package dtos

import (
	"encoding/json"
	"testing"
)

// decode mimics what the Bambu client does: unmarshal the `print` object of a
// report onto the accumulated state.
func decode(t *testing.T, state *BambuPrintState, raw string) *BambuPrintState {
	t.Helper()
	var envelope struct {
		Print json.RawMessage `json:"print"`
	}
	if err := json.Unmarshal([]byte(raw), &envelope); err != nil {
		t.Fatalf("envelope decode: %v", err)
	}
	if state == nil {
		state = &BambuPrintState{}
	}
	if err := json.Unmarshal(envelope.Print, state); err != nil {
		t.Fatalf("state decode: %v", err)
	}
	return state
}

func TestCoercionAndProjection(t *testing.T) {
	// nozzle_temp_min/max arrive as numeric strings, temps as numbers.
	full := `{"print":{
		"gcode_state":"RUNNING",
		"stg_cur":2,
		"mc_percent":42,
		"mc_remaining_time":17,
		"file":"benchy.3mf",
		"bed_temper":59.5,
		"bed_target_temper":60,
		"nozzle_temper":219,
		"nozzle_target_temper":220,
		"support_chamber_temp_edit":true,
		"device":{"ctc":{"info":{"temp":2949180}}},
		"ams":{"tray_now":"1","ams":[{"id":0,"humidity_raw":"35","dry_time":0,"temp":"28.5","tray":[
			{"id":0,"tray_type":"PLA","tray_sub_brands":"Bambu","tray_color":"FF0000","nozzle_temp_min":"190","nozzle_temp_max":"230","remain":75},
			{"id":1,"tray_type":""}
		]}]},
		"vir_slot":[{"id":254,"tray_type":"PETG"}],
		"lights_report":[{"node":"chamber_light","mode":"on"}]
	}}`

	state := decode(t, nil, full)
	s := StatusFromMQTT(state)

	if s.State != "RUNNING" || s.Progress != 42 {
		t.Fatalf("basic fields: %+v", s)
	}
	if s.Stage == nil || *s.Stage != 2 {
		t.Fatalf("stage: %+v", s.Stage)
	}
	if s.TimeRemaining == nil || *s.TimeRemaining != 17 {
		t.Fatalf("timeRemaining: %+v", s.TimeRemaining)
	}
	// chamber: 2949180 = (45 << 16) | 60
	if s.Chamber.Temperature != 60 || s.Chamber.TargetTemperature != 45 || !s.Chamber.Controllable {
		t.Fatalf("chamber: %+v", s.Chamber)
	}
	if len(s.AMS) != 1 || s.AMS[0].Humidity != 35 || s.AMS[0].Temperature != 28.5 {
		t.Fatalf("ams unit: %+v", s.AMS)
	}
	// tray_now is "1" → globalId 1 (unit 0, tray 1) is the loaded tray.
	tray0 := s.AMS[0].Trays[0]
	if tray0.Empty || tray0.Loaded || tray0.Material == nil || *tray0.Material != "PLA" {
		t.Fatalf("tray0: %+v", tray0)
	}
	if tray0.NozzleTempMin == nil || *tray0.NozzleTempMin != 190 || tray0.NozzleTempMax == nil || *tray0.NozzleTempMax != 230 {
		t.Fatalf("tray0 temps: %+v", tray0)
	}
	tray1 := s.AMS[0].Trays[1]
	if !tray1.Empty || !tray1.Loaded || tray1.Material != nil {
		t.Fatalf("tray1: %+v", tray1)
	}
	if s.ExternalSpool == nil || s.ExternalSpool.ID != 254 || s.ExternalSpool.Material == nil || *s.ExternalSpool.Material != "PETG" {
		t.Fatalf("externalSpool: %+v", s.ExternalSpool)
	}
	if !s.ChamberLight {
		t.Fatalf("chamberLight should be on")
	}
}

func TestFanAndSpeedProjection(t *testing.T) {
	// Fan gears arrive as numeric strings on a 0-15 scale; spd_lvl as a number.
	state := decode(t, nil, `{"print":{
		"gcode_state":"RUNNING",
		"spd_lvl":3,
		"cooling_fan_speed":"15",
		"big_fan1_speed":"8",
		"big_fan2_speed":"0"
	}}`)

	s := StatusFromMQTT(state)

	if s.SpeedLevel == nil || *s.SpeedLevel != 3 {
		t.Fatalf("speedLevel: %+v", s.SpeedLevel)
	}
	// 15/15 -> 100, 8/15 -> 53, 0 -> 0
	if s.Fans.Part != 100 || s.Fans.Aux != 53 || s.Fans.Chamber != 0 {
		t.Fatalf("fans: %+v", s.Fans)
	}
}

func TestSpeedLevelOmittedWhenAbsent(t *testing.T) {
	state := decode(t, nil, `{"print":{"gcode_state":"IDLE"}}`)
	s := StatusFromMQTT(state)
	if s.SpeedLevel != nil {
		t.Fatalf("speedLevel should be nil when absent: %+v", s.SpeedLevel)
	}
	if s.Fans.Part != 0 || s.Fans.Aux != 0 || s.Fans.Chamber != 0 {
		t.Fatalf("fans should default to zero: %+v", s.Fans)
	}
}

func TestPartialMergePreservesAndReplaces(t *testing.T) {
	state := decode(t, nil, `{"print":{
		"gcode_state":"RUNNING",
		"nozzle_temper":200,
		"ams":{"tray_now":"0","ams":[{"id":0,"tray":[{"id":0},{"id":1}]}]},
		"lights_report":[{"node":"chamber_light","mode":"on"}]
	}}`)

	// A later partial report only updates a couple of scalars and the ams array.
	state = decode(t, state, `{"print":{
		"mc_percent":10,
		"ams":{"ams":[{"id":0,"tray":[{"id":0,"tray_type":"PLA"}]}]}
	}}`)

	s := StatusFromMQTT(state)

	// Preserved from the first message.
	if s.State != "RUNNING" || s.Nozzle.Temperature != 200 {
		t.Fatalf("scalars not preserved: %+v", s)
	}
	if !s.ChamberLight {
		t.Fatalf("lights_report not preserved")
	}
	// tray_now was not in the second message, so it is preserved ("0").
	// The ams array was replaced wholesale: now a single tray.
	if len(s.AMS) != 1 || len(s.AMS[0].Trays) != 1 {
		t.Fatalf("ams array not replaced: %+v", s.AMS)
	}
	if s.Progress != 10 {
		t.Fatalf("progress not updated: %v", s.Progress)
	}
}
