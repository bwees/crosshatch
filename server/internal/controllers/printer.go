package controllers

import (
	"crosshatch/internal/database/models"
	"crosshatch/internal/dtos"
	"crosshatch/internal/services"
	"strconv"

	"github.com/go-fuego/fuego"
)

type PrinterController struct {
	svc *services.PrinterService
}

func (c *PrinterController) Register(api *fuego.Server) {
	route := fuego.Group(api, "/printer")

	fuego.Get(route, "/", func(ctx fuego.ContextNoBody) ([]models.Printer, error) {
		return c.svc.GetPrinters()
	},
		fuego.OptionOperationID("getPrinters"),
	)

	fuego.Get(route, "/{serial}/status", func(ctx fuego.ContextNoBody) (*dtos.PrinterStatus, error) {
		return c.svc.GetStatus(ctx.PathParam("serial"))
	},
		fuego.OptionOperationID("getPrinterStatus"),
	)

	fuego.Put(route, "/", func(ctx fuego.ContextWithBody[dtos.CreatePrinterDto]) (*models.Printer, error) {
		dto, err := ctx.Body()
		if err != nil {
			return nil, err
		}

		return c.svc.CreatePrinter(dto)
	},
		fuego.OptionOperationID("createPrinter"),
		fuego.OptionDefaultStatusCode(201),
	)

	fuego.Patch(route, "/{serial}", func(ctx fuego.ContextWithBody[dtos.UpdatePrinterDto]) (*models.Printer, error) {
		dto, err := ctx.Body()
		if err != nil {
			return nil, err
		}

		return c.svc.UpdatePrinter(ctx.PathParam("serial"), dto)
	},
		fuego.OptionOperationID("updatePrinter"),
	)

	fuego.Delete(route, "/{serial}", func(ctx fuego.ContextNoBody) (any, error) {
		return nil, c.svc.DeletePrinter(ctx.PathParam("serial"))
	},
		fuego.OptionOperationID("deletePrinter"),
		fuego.OptionDefaultStatusCode(204),
	)

	fuego.Post(route, "/{serial}/stop", func(ctx fuego.ContextNoBody) (any, error) {
		return nil, c.svc.StopPrint(ctx.PathParam("serial"))
	},
		fuego.OptionOperationID("stopPrint"),
		fuego.OptionDefaultStatusCode(204),
	)

	fuego.Post(route, "/{serial}/pause", func(ctx fuego.ContextNoBody) (any, error) {
		return nil, c.svc.PausePrint(ctx.PathParam("serial"))
	},
		fuego.OptionOperationID("pausePrint"),
		fuego.OptionDefaultStatusCode(204),
	)

	fuego.Post(route, "/{serial}/resume", func(ctx fuego.ContextNoBody) (any, error) {
		return nil, c.svc.ResumePrint(ctx.PathParam("serial"))
	},
		fuego.OptionOperationID("resumePrint"),
		fuego.OptionDefaultStatusCode(204),
	)

	fuego.Post(route, "/{serial}/light", func(ctx fuego.ContextWithBody[dtos.SetLightDto]) (any, error) {
		dto, err := ctx.Body()
		if err != nil {
			return nil, err
		}
		return nil, c.svc.SetLight(ctx.PathParam("serial"), dto.State)
	},
		fuego.OptionOperationID("setLight"),
		fuego.OptionDefaultStatusCode(204),
	)

	fuego.Post(route, "/{serial}/speed", func(ctx fuego.ContextWithBody[dtos.SetPrintSpeedDto]) (any, error) {
		dto, err := ctx.Body()
		if err != nil {
			return nil, err
		}
		return nil, c.svc.SetPrintSpeed(ctx.PathParam("serial"), dto.Level)
	},
		fuego.OptionOperationID("setPrintSpeed"),
		fuego.OptionDefaultStatusCode(204),
	)

	fuego.Post(route, "/{serial}/unload/{amsId}", func(ctx fuego.ContextNoBody) (any, error) {
		amsID, err := strconv.Atoi(ctx.PathParam("amsId"))
		if err != nil {
			return nil, err
		}
		return nil, c.svc.UnloadMaterial(ctx.PathParam("serial"), amsID)
	},
		fuego.OptionOperationID("unloadMaterial"),
		fuego.OptionDefaultStatusCode(204),
	)
}

func NewPrinterController(svc *services.PrinterService) *PrinterController {
	return &PrinterController{svc: svc}
}
