package main

import (
	"crosshatch/internal/controllers"
	"crosshatch/internal/database"
	"crosshatch/internal/repositories"
	"crosshatch/internal/services"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		// fx.NopLogger,
		controllers.Module,
		services.Module,
		database.Module,
		repositories.Module,

		fx.Invoke(NewServer),
	)

	app.Run()
}
