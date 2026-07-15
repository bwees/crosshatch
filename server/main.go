package main

import (
	"crosshatch/internal/controllers"
	"crosshatch/internal/database"
	"crosshatch/internal/go2rtc"
	"crosshatch/internal/proxy"
	"crosshatch/internal/repositories"
	"crosshatch/internal/services"
	"crosshatch/internal/socketio"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.NopLogger,
		controllers.Module,
		services.Module,
		database.Module,
		repositories.Module,
		socketio.Module,
		proxy.Module,
		go2rtc.Module,

		fx.Invoke(NewServer),
	)

	app.Run()
}
