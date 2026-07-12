package services

import "go.uber.org/fx"

var Module = fx.Provide(
	NewPrinterService,
	NewFilamentService,
	NewAuthService,
)
