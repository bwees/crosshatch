package controllers

import (
	"crosshatch/internal/database/models"
	"crosshatch/internal/dtos"
	"crosshatch/internal/repositories"
	"crosshatch/internal/services"

	"github.com/go-fuego/fuego"
)

type NotificationController struct {
	vapid   *services.VapidService
	repo    *repositories.NotificationRepository
	service *services.NotificationService
}

func (c *NotificationController) Register(api *fuego.Server) {
	notifications := fuego.Group(api, "/notifications")

	fuego.Get(notifications, "/vapid", func(ctx fuego.ContextNoBody) (dtos.VapidDto, error) {
		return c.vapid.GetConfig()
	},
		fuego.OptionOperationID("getVapidConfig"),
	)

	fuego.Post(notifications, "/subscribe", func(ctx fuego.ContextWithBody[dtos.PushSubscriptionDto]) (any, error) {
		user := userFromContext(ctx.Request().Context())
		if user == nil {
			return nil, fuego.UnauthorizedError{Title: "Unauthorized"}
		}

		dto, err := ctx.Body()
		if err != nil {
			return nil, err
		}

		return nil, c.repo.UpsertSubscription(dto.DeviceId, user.ID, dto.Endpoint, dto.P256dh, dto.Auth)
	},
		fuego.OptionOperationID("subscribeNotifications"),
		fuego.OptionDefaultStatusCode(204),
	)

	fuego.Post(notifications, "/unsubscribe", func(ctx fuego.ContextWithBody[dtos.UnsubscribeDto]) (any, error) {
		user := userFromContext(ctx.Request().Context())
		if user == nil {
			return nil, fuego.UnauthorizedError{Title: "Unauthorized"}
		}

		dto, err := ctx.Body()
		if err != nil {
			return nil, err
		}

		return nil, c.repo.DeleteSubscriptionByDevice(dto.DeviceId)
	},
		fuego.OptionOperationID("unsubscribeNotifications"),
		fuego.OptionDefaultStatusCode(204),
	)

	fuego.Post(notifications, "/test", func(ctx fuego.ContextWithBody[dtos.TestNotificationDto]) (any, error) {
		user := userFromContext(ctx.Request().Context())
		if user == nil {
			return nil, fuego.UnauthorizedError{Title: "Unauthorized"}
		}

		dto, err := ctx.Body()
		if err != nil {
			return nil, err
		}

		return nil, c.service.SendTest(dto.DeviceId, dto.Serial)
	},
		fuego.OptionOperationID("testNotification"),
		fuego.OptionDefaultStatusCode(204),
	)

	fuego.Get(notifications, "/settings/{deviceId}/{serial}", func(ctx fuego.ContextNoBody) (dtos.NotificationSettingsDto, error) {
		user := userFromContext(ctx.Request().Context())
		if user == nil {
			return dtos.NotificationSettingsDto{}, fuego.UnauthorizedError{Title: "Unauthorized"}
		}

		setting, err := c.repo.GetSetting(ctx.PathParam("deviceId"), ctx.PathParam("serial"))
		if err != nil {
			return dtos.NotificationSettingsDto{}, err
		}
		return setting.ToDto(), nil
	},
		fuego.OptionOperationID("getPrinterNotificationSettings"),
	)

	fuego.Put(notifications, "/settings/{deviceId}/{serial}", func(ctx fuego.ContextWithBody[dtos.NotificationSettingsDto]) (dtos.NotificationSettingsDto, error) {
		user := userFromContext(ctx.Request().Context())
		if user == nil {
			return dtos.NotificationSettingsDto{}, fuego.UnauthorizedError{Title: "Unauthorized"}
		}

		dto, err := ctx.Body()
		if err != nil {
			return dtos.NotificationSettingsDto{}, err
		}

		setting := models.NotificationSetting{
			DeviceID:       ctx.PathParam("deviceId"),
			PrinterSerial:  ctx.PathParam("serial"),
			UserID:         user.ID,
			Enabled:        dto.Enabled,
			NotifyComplete: dto.NotifyComplete,
			NotifyError:    dto.NotifyError,
		}
		if err := c.repo.UpsertSetting(setting); err != nil {
			return dtos.NotificationSettingsDto{}, err
		}
		return setting.ToDto(), nil
	},
		fuego.OptionOperationID("updatePrinterNotificationSettings"),
	)
}

func NewNotificationController(vapid *services.VapidService, repo *repositories.NotificationRepository, service *services.NotificationService) *NotificationController {
	return &NotificationController{vapid: vapid, repo: repo, service: service}
}
