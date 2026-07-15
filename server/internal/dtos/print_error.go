package dtos

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
)

// printErrorsRaw is the bundled Bambu device error map, keyed by the uppercase
// 8-hex-char code, ported from OrcaSlicer's resources/hms device_error tables.
//
//go:embed print_errors.json
var printErrorsRaw []byte

var printErrorMessages = loadPrintErrors()

func loadPrintErrors() map[string]string {
	var m map[string]string
	if err := json.Unmarshal(printErrorsRaw, &m); err != nil {
		slog.Error("failed to load print error map", "error", err)
		return map[string]string{}
	}
	return m
}

// cancelCodes are the device error codes reported when a print is cancelled by
// the user, as opposed to failing on a fault.
var cancelCodes = map[uint32]bool{
	0x0300400C: true, // The task was canceled.
	0x0500400E: true, // Printing was cancelled.
}

// PrinterError is a decoded Bambu print_error code.
type PrinterError struct {
	Code      string `json:"code" validate:"required"`    // dashed for display, e.g. "0300-400C"
	Summary   string `json:"summary" validate:"required"` // first sentence, for tight spaces
	Message   string `json:"message,omitempty"`           // full text, empty for unknown codes
	Cancelled bool   `json:"cancelled"`                   // user cancel rather than a fault
}

// DecodePrintError turns a raw Bambu print_error value into a PrinterError, or
// nil when the code is zero (no active error). Unknown codes still yield an
// error carrying the formatted code with an empty message.
func DecodePrintError(raw uint32) *PrinterError {
	if raw == 0 {
		return nil
	}

	hex := fmt.Sprintf("%08X", raw)
	code := hex[:4] + "-" + hex[4:]
	message := printErrorMessages[hex]

	summary := firstSentence(message)
	if summary == "" {
		summary = "Error " + code
	}

	return &PrinterError{
		Code:      code,
		Summary:   summary,
		Message:   message,
		Cancelled: cancelCodes[raw],
	}
}

// firstSentence returns the leading sentence of a message. Bambu's messages are
// structured as "<problem>. <how to fix>.", so the first sentence is a compact
// form for notification bodies and status cards.
func firstSentence(msg string) string {
	if i := strings.Index(msg, ". "); i != -1 {
		return msg[:i+1]
	}
	return msg
}
