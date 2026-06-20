package controllers

import (
	"crosshatch/internal/database/models"
	"crosshatch/internal/services"

	"github.com/go-fuego/fuego"
)

type PrinterController struct {
	svc *services.PrinterService
}

func (c *PrinterController) Register(api *fuego.Server) {
	route := fuego.Group(api, "/printer")

	fuego.Get(route, "/", func(ctx fuego.ContextNoBody) ([]models.Printer, error) {
		return c.svc.GetPrinters()
	})

}

func NewPrinterController(svc *services.PrinterService) *PrinterController {
	return &PrinterController{svc: svc}
}
