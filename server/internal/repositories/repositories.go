package repositories

import "go.uber.org/fx"

var Module = fx.Provide(
	NewPrinterRepository,
	NewCameraRepository,
	NewFilamentRepository,
	NewUserRepository,
	NewSessionRepository,
)
