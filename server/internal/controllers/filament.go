package controllers

import (
	"crosshatch/internal/database/models"
	"crosshatch/internal/services"

	"github.com/go-fuego/fuego"
)

type FilamentController struct {
	svc *services.FilamentService
}

func (c *FilamentController) Register(api *fuego.Server) {
	route := fuego.Group(api, "/filament")

	fuego.Get(route, "/", func(ctx fuego.ContextNoBody) ([]models.Filament, error) {
		return c.svc.GetFilaments()
	},
		fuego.OptionOperationID("getFilaments"),
	)
}

func NewFilamentController(svc *services.FilamentService) *FilamentController {
	return &FilamentController{svc: svc}
}
