package dtos

import "testing"

func TestDecodePrintError(t *testing.T) {
	if DecodePrintError(0) != nil {
		t.Fatal("code 0 should decode to nil")
	}

	cancel := DecodePrintError(0x0300400C)
	if cancel == nil || !cancel.Cancelled || cancel.Code != "0300-400C" {
		t.Fatalf("cancel decode: %+v", cancel)
	}

	// Known fault: full message present, summary is the first sentence only.
	clog := DecodePrintError(0x03008016)
	if clog == nil || clog.Cancelled {
		t.Fatalf("clog decode: %+v", clog)
	}
	if clog.Summary != "The nozzle is clogged with filament." {
		t.Fatalf("clog summary: %q", clog.Summary)
	}
	if len(clog.Message) <= len(clog.Summary) {
		t.Fatalf("full message should exceed summary: %q", clog.Message)
	}

	// Unknown code still yields a usable error keyed by its formatted code.
	unknown := DecodePrintError(0xDEADBEEF)
	if unknown == nil || unknown.Code != "DEAD-BEEF" || unknown.Message != "" {
		t.Fatalf("unknown decode: %+v", unknown)
	}
	if unknown.Summary != "Error DEAD-BEEF" {
		t.Fatalf("unknown summary: %q", unknown.Summary)
	}
}

func TestPrintErrorProjection(t *testing.T) {
	state := decode(t, nil, `{"print":{"gcode_state":"FAILED","print_error":50348044}}`)
	s := StatusFromMQTT(state)
	if s.Error == nil || !s.Error.Cancelled {
		t.Fatalf("expected cancelled error, got %+v", s.Error)
	}

	clear := decode(t, state, `{"print":{"print_error":0}}`)
	if StatusFromMQTT(clear).Error != nil {
		t.Fatal("print_error 0 should clear the error")
	}
}
