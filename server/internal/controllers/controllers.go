package controllers

import (
	"github.com/go-fuego/fuego"
	"go.uber.org/fx"
)

type Controller interface {
	Register(*fuego.Server)
}

var Module = fx.Provide(
	fx.Annotate(NewPrinterController, fx.As(new(Controller)), fx.ResultTags(`group:"controllers"`)),
	fx.Annotate(NewWebsocketController, fx.As(new(Controller)), fx.ResultTags(`group:"controllers"`)),
	fx.Annotate(NewFilamentController, fx.As(new(Controller)), fx.ResultTags(`group:"controllers"`)),
	fx.Annotate(NewAuthController, fx.As(new(Controller)), fx.ResultTags(`group:"controllers"`)),
	fx.Annotate(NewUsersController, fx.As(new(Controller)), fx.ResultTags(`group:"controllers"`)),
	fx.Annotate(NewNotificationController, fx.As(new(Controller)), fx.ResultTags(`group:"controllers"`)),
	NewAuthMiddleware,
)
