package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"crosshatch/internal/controllers"

	"github.com/go-fuego/fuego"
	"go.uber.org/fx"
)

type Controllers struct {
	fx.In
	All []controllers.Controller `group:"controllers"`
}

func NewServer(lc fx.Lifecycle, controllers Controllers) *fuego.Server {
	server := fuego.NewServer(
		fuego.WithLoggingMiddleware(fuego.LoggingConfig{
			DisableRequest:  true,
			DisableResponse: true,
		}),
		fuego.WithSerializer(fuego.SendJSON),
	)

	server.OpenAPI.Config.JSONFilePath = "openapi.json"
	server.OpenAPI.Config.SpecURL = "/openapi.json"
	server.OpenAPI.Config.DisableSwaggerUI = true

	api := fuego.Group(server, "/api")
	for _, controller := range controllers.All {
		controller.Register(api)
		slog.Info(fmt.Sprintf("Registered controller: %T\n", controller))
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			openapi := flag.Bool("open-api", false, "Generate OpenAPI documentation")
			flag.Parse()

			if *openapi {
				server.OutputOpenAPISpec()
				slog.Info("OpenAPI documentation generated successfully.")
				os.Exit(0)
			}

			go server.Run()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.Shutdown(ctx)
			return nil
		},
	})

	return server
}
