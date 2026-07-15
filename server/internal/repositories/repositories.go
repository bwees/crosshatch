package repositories

import "go.uber.org/fx"

var Module = fx.Provide(
	NewPrinterRepository,
	NewFilamentRepository,
	NewUserRepository,
	NewSessionRepository,
	NewNotificationRepository,
)
