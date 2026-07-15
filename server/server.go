package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"

	"crosshatch/internal/config"
	"crosshatch/internal/controllers"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	"go.uber.org/fx"
)

type Controllers struct {
	fx.In
	All []controllers.Controller `group:"controllers"`
}

func NewServer(lc fx.Lifecycle, controllers Controllers, authMiddleware *controllers.AuthMiddleware) *fuego.Server {
	addr := ":" + config.Port()

	server := fuego.NewServer(
		fuego.WithLoggingMiddleware(fuego.LoggingConfig{
			DisableRequest:  true,
			DisableResponse: true,
		}),
		fuego.WithSerializer(fuego.SendJSON),
		fuego.WithAddr(addr),
	)

	server.OpenAPI.Config.JSONFilePath = "openapi.json"
	server.OpenAPI.Config.SpecURL = "/openapi.json"
	server.OpenAPI.Config.DisableSwaggerUI = true
	server.OpenAPI.Config.PrettyFormatJSON = true

	// Optional fields use Go pointers with `omitempty`: when nil they are
	// omitted from the JSON entirely, never serialized as null. Fuego marks
	// pointer fields nullable by default, which makes the generated client
	// types `T | null`. Override that so they generate as `T | undefined`,
	// matching how the values actually appear on the wire.
	// Fuego runs its own default customizer before this one, so we only apply
	// the nullable override here (calling the default again would re-run its
	// required-field detection after the markers are consumed, wiping it).
	server.OpenAPI.SetGeneratorSchemaCustomizer(func(name string, t reflect.Type, tag reflect.StructTag, schema *openapi3.Schema) error {
		if strings.Contains(tag.Get("json"), "omitempty") {
			schema.Nullable = false
		}
		return nil
	})

	fuego.Use(server, authMiddleware.Handler)

	api := fuego.Group(server, "/api")
	for _, controller := range controllers.All {
		controller.Register(api)
		slog.Info(fmt.Sprintf("Registered controller: %T\n", controller))
	}

	// Serve the built web frontend (SvelteKit static build) when configured.
	// In production the Go server is the single origin for both the API and
	// the app; in dev the web is served by Vite instead, so WEB_STATIC_PATH is
	// left unset and this is a no-op.
	registerStaticWeb(server)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			openapi := flag.Bool("open-api", false, "Generate OpenAPI documentation")
			flag.Parse()

			if *openapi {
				server.OutputOpenAPISpec()
				slog.Info("OpenAPI documentation generated successfully.")
				os.Exit(0)
			}

			go func() {
				if err := server.Run(); err != nil {
					slog.Error("Failed to start server", "error", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.Shutdown(ctx)
			return nil
		},
	})

	return server
}
