package services

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewPrinterService,
		NewFilamentService,
		NewAuthService,
		NewVapidService,
		NewNotificationService,
	),
	// NotificationService self-registers as a printer status observer with side
	// effects, so it must be constructed even though nothing depends on it.
	fx.Invoke(func(*NotificationService) {}),
)
