package services

import (
	"testing"

	"crosshatch/internal/dtos"
)

func TestClassifyTransition(t *testing.T) {
	status := func(state string) *dtos.PrinterStatus {
		return &dtos.PrinterStatus{State: state}
	}

	cases := []struct {
		name      string
		prev      *dtos.PrinterStatus
		next      *dtos.PrinterStatus
		wantEvent string
		wantOk    bool
	}{
		{"running to finish is complete", status("RUNNING"), status("FINISH"), "complete", true},
		{"running to failed is error", status("RUNNING"), status("FAILED"), "error", true},
		{"failed to failed is none", status("FAILED"), status("FAILED"), "", false},
		{"running to pause is none", status("RUNNING"), status("PAUSE"), "", false},
		{"nil prev to finish is none", nil, status("FINISH"), "", false},
		{"nil prev to failed is error", nil, status("FAILED"), "error", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			event, ok := classifyTransition(tc.prev, tc.next)
			if event != tc.wantEvent || ok != tc.wantOk {
				t.Fatalf("classifyTransition() = (%q, %v), want (%q, %v)", event, ok, tc.wantEvent, tc.wantOk)
			}
		})
	}
}
