package services

import (
	"testing"

	"crosshatch/internal/dtos"
)

func TestClassifyTransition(t *testing.T) {
	status := func(state string) *dtos.PrinterStatus {
		return &dtos.PrinterStatus{State: state}
	}
	withError := func(state string, err *dtos.PrinterError) *dtos.PrinterStatus {
		return &dtos.PrinterStatus{State: state, Error: err}
	}

	cancel := &dtos.PrinterError{Code: "0300-400C", Summary: "The task was canceled.", Cancelled: true}
	fault := &dtos.PrinterError{Code: "0300-8016", Summary: "The nozzle is clogged with filament."}

	cases := []struct {
		name      string
		prev      *dtos.PrinterStatus
		next      *dtos.PrinterStatus
		wantEvent dtos.NotificationEvent
		wantTitle string
		wantOk    bool
	}{
		{"running to finish is complete", status("RUNNING"), status("FINISH"), dtos.EventComplete, "Print complete", true},
		{"running to failed is error", status("RUNNING"), status("FAILED"), dtos.EventError, "Print error", true},
		{"failed to failed is none", status("FAILED"), status("FAILED"), "", "", false},
		{"running to pause is none", status("RUNNING"), status("PAUSE"), "", "", false},
		{"nil prev to finish is none", nil, status("FINISH"), "", "", false},
		{"nil prev to failed is error", nil, status("FAILED"), dtos.EventError, "Print error", true},
		{"new cancel code is cancelled", status("RUNNING"), withError("FAILED", cancel), dtos.EventError, "Print cancelled", true},
		{"new fault code is error", status("RUNNING"), withError("PAUSE", fault), dtos.EventError, "Print error", true},
		{"persisting error does not re-fire", withError("FAILED", fault), withError("FAILED", fault), "", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg, ok := classifyTransition("Printer", tc.prev, tc.next)
			if msg.Event != tc.wantEvent || ok != tc.wantOk || msg.Title != tc.wantTitle {
				t.Fatalf("classifyTransition() = (%q, %q, %v), want (%q, %q, %v)", msg.Event, msg.Title, ok, tc.wantEvent, tc.wantTitle, tc.wantOk)
			}
		})
	}
}
