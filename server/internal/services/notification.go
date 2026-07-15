package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"crosshatch/internal/dtos"
	"crosshatch/internal/repositories"

	webpush "github.com/SherClockHolmes/webpush-go"
)

type NotificationService struct {
	printers      *repositories.PrinterRepository
	notifications *repositories.NotificationRepository
	vapid         *VapidService
}

type pushPayload struct {
	PrinterSerial string `json:"printerSerial"`
	Title         string `json:"title"`
	Body          string `json:"body"`
	Tag           string `json:"tag"`
}

// pushMessage is a notification ready to deliver: the event drives per-printer
// subscription filtering, while title/body are the rendered strings.
type pushMessage struct {
	Event dtos.NotificationEvent
	Title string
	Body  string
}

// classifyTransition maps a status change to a notification for the printer
// named `name`. The second return is true when the transition warrants one.
func classifyTransition(name string, prev, next *dtos.PrinterStatus) (pushMessage, bool) {
	if next == nil {
		return pushMessage{}, false
	}

	// A newly-reported print error takes precedence and fires regardless of
	// gcode_state, so mid-print faults that only pause are still surfaced.
	if next.Error != nil && (prev == nil || prev.Error == nil) {
		if next.Error.Cancelled {
			return pushMessage{dtos.EventError, "Print cancelled", name + " print was cancelled"}, true
		}
		return pushMessage{dtos.EventError, "Print error", name + ": " + next.Error.Summary}, true
	}

	if prev != nil && prev.State == dtos.GcodeRunning && next.State == dtos.GcodeFinish {
		return pushMessage{dtos.EventComplete, "Print complete", name + " finished printing"}, true
	}

	// Fallback: entered FAILED without a decoded error code.
	if next.State == dtos.GcodeFailed && (prev == nil || prev.State != dtos.GcodeFailed) {
		return pushMessage{dtos.EventError, "Print error", name + " reported an error"}, true
	}

	return pushMessage{}, false
}

func (s *NotificationService) onTransition(serial string, prev, next *dtos.PrinterStatus) {
	msg, ok := classifyTransition(s.printerName(serial), prev, next)
	if !ok {
		return
	}
	s.deliver(serial, msg)
}

// printerName returns the printer's display name, falling back to its serial.
func (s *NotificationService) printerName(serial string) string {
	if printer, err := s.printers.GetPrinterBySerial(serial); err == nil && printer != nil {
		return printer.Name
	}
	return serial
}

// deliver sends a push notification to every device subscribed to its event
// for this printer.
func (s *NotificationService) deliver(serial string, msg pushMessage) {
	body, err := json.Marshal(pushPayload{
		PrinterSerial: serial,
		Title:         msg.Title,
		Body:          msg.Body,
		Tag:           "crosshatch-" + serial,
	})
	if err != nil {
		slog.Error("failed to marshal notification payload", "error", err)
		return
	}

	subs, err := s.notifications.SubscriptionsForEvent(serial, msg.Event)
	if err != nil {
		slog.Error("failed to find subscriptions", "serial", serial, "error", err)
		return
	}

	for _, sub := range subs {
		s.send(sub.Endpoint, sub.P256dh, sub.Auth, body)
	}
}

// SendTest delivers a test notification directly to a single device, bypassing
// the per-printer settings so a user can verify push works on this device.
func (s *NotificationService) SendTest(deviceID, serial string) error {
	sub, err := s.notifications.GetSubscriptionByDevice(deviceID)
	if err != nil {
		return err
	}
	if sub == nil {
		return fmt.Errorf("no push subscription for this device")
	}

	body, err := json.Marshal(pushPayload{
		PrinterSerial: serial,
		Title:         "Test notification",
		Body:          s.printerName(serial) + " notifications are working",
		Tag:           "crosshatch-" + serial,
	})
	if err != nil {
		return err
	}

	s.send(sub.Endpoint, sub.P256dh, sub.Auth, body)
	return nil
}

func (s *NotificationService) send(endpoint, p256dh, auth string, body []byte) {
	sub := &webpush.Subscription{
		Endpoint: endpoint,
		Keys:     webpush.Keys{P256dh: p256dh, Auth: auth},
	}

	res, err := webpush.SendNotification(body, sub, &webpush.Options{
		Subscriber:      s.vapid.Subject(),
		VAPIDPublicKey:  s.vapid.PublicKey(),
		VAPIDPrivateKey: s.vapid.PrivateKey(),
		TTL:             60,
	})
	if err != nil {
		slog.Error("failed to send web push", "endpoint", endpoint, "error", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound || res.StatusCode == http.StatusGone {
		if err := s.notifications.DeleteSubscriptionByEndpoint(endpoint); err != nil {
			slog.Error("failed to delete stale subscription", "endpoint", endpoint, "error", err)
		}
		return
	}

	if res.StatusCode >= 300 {
		responseBody, _ := io.ReadAll(res.Body)
		slog.Warn("web push rejected", "endpoint", endpoint, "status", res.StatusCode, "body", string(responseBody))
	}
}

func NewNotificationService(
	printerSvc *PrinterService,
	printers *repositories.PrinterRepository,
	notifications *repositories.NotificationRepository,
	vapid *VapidService,
) *NotificationService {
	svc := &NotificationService{
		printers:      printers,
		notifications: notifications,
		vapid:         vapid,
	}

	printerSvc.AddObserver(svc.onTransition)

	return svc
}
