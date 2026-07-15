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

// classifyTransition maps a status change to a notification event. It returns
// ("complete"|"error", true) when the transition warrants a notification, and
// ("", false) otherwise.
func classifyTransition(prev, next *dtos.PrinterStatus) (string, bool) {
	if next == nil {
		return "", false
	}

	if prev != nil && prev.State == "RUNNING" && next.State == "FINISH" {
		return "complete", true
	}

	if next.State == "FAILED" && (prev == nil || prev.State != "FAILED") {
		return "error", true
	}

	return "", false
}

func (s *NotificationService) onTransition(serial string, prev, next *dtos.PrinterStatus) {
	event, ok := classifyTransition(prev, next)
	if !ok {
		return
	}
	s.Notify(serial, event)
}

// printerName returns the printer's display name, falling back to its serial.
func (s *NotificationService) printerName(serial string) string {
	if printer, err := s.printers.GetPrinterBySerial(serial); err == nil && printer != nil {
		return printer.Name
	}
	return serial
}

// Notify builds and delivers a push notification for the given event
// ("complete" or "error") to every device subscribed to it for this printer.
func (s *NotificationService) Notify(serial, event string) {
	name := s.printerName(serial)

	payload := pushPayload{
		PrinterSerial: serial,
		Tag:           "crosshatch-" + serial,
	}
	switch event {
	case "complete":
		payload.Title = "Print complete"
		payload.Body = name + " finished printing"
	case "error":
		payload.Title = "Print error"
		payload.Body = name + " reported an error"
	}

	body, err := json.Marshal(payload)
	if err != nil {
		slog.Error("failed to marshal notification payload", "error", err)
		return
	}

	subs, err := s.notifications.SubscriptionsForEvent(serial, event)
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
