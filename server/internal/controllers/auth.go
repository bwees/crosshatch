package controllers

import (
	"crosshatch/internal/database/models"
	"crosshatch/internal/dtos"
	"crosshatch/internal/services"

	"github.com/go-fuego/fuego"
)

type AuthController struct {
	svc *services.AuthService
}

func userDto(u *models.User) dtos.UserDto {
	return dtos.UserDto{ID: u.ID, Username: u.Username, IsAdmin: u.IsAdmin}
}

func (c *AuthController) Register(api *fuego.Server) {
	route := fuego.Group(api, "/auth")

	fuego.Get(route, "/setup", func(ctx fuego.ContextNoBody) (dtos.SetupStatusDto, error) {
		required, err := c.svc.SetupRequired()
		if err != nil {
			return dtos.SetupStatusDto{}, err
		}
		return dtos.SetupStatusDto{SetupRequired: required}, nil
	},
		fuego.OptionOperationID("getSetupStatus"),
	)

	fuego.Post(route, "/setup", func(ctx fuego.ContextWithBody[dtos.CreateUserDto]) (dtos.UserDto, error) {
		dto, err := ctx.Body()
		if err != nil {
			return dtos.UserDto{}, err
		}

		user, err := c.svc.Setup(dto)
		if err != nil {
			return dtos.UserDto{}, err
		}

		_, token, err := c.svc.Login(dtos.LoginDto{Username: dto.Username, Password: dto.Password})
		if err != nil {
			return dtos.UserDto{}, err
		}

		secure := ctx.Request().TLS != nil
		ctx.SetCookie(sessionCookie(token, secure, int(services.SessionDuration.Seconds())))
		return userDto(user), nil
	},
		fuego.OptionOperationID("setup"),
	)

	fuego.Post(route, "/login", func(ctx fuego.ContextWithBody[dtos.LoginDto]) (dtos.UserDto, error) {
		dto, err := ctx.Body()
		if err != nil {
			return dtos.UserDto{}, err
		}

		user, token, err := c.svc.Login(dto)
		if err != nil {
			return dtos.UserDto{}, err
		}

		secure := ctx.Request().TLS != nil
		ctx.SetCookie(sessionCookie(token, secure, int(services.SessionDuration.Seconds())))
		return userDto(user), nil
	},
		fuego.OptionOperationID("login"),
	)

	fuego.Post(route, "/logout", func(ctx fuego.ContextNoBody) (any, error) {
		if cookie, err := ctx.Cookie(sessionCookieName); err == nil {
			if err := c.svc.Logout(cookie.Value); err != nil {
				return nil, err
			}
		}

		secure := ctx.Request().TLS != nil
		ctx.SetCookie(sessionCookie("", secure, -1))
		return nil, nil
	},
		fuego.OptionOperationID("logout"),
		fuego.OptionDefaultStatusCode(204),
	)

	fuego.Get(route, "/me", func(ctx fuego.ContextNoBody) (dtos.UserDto, error) {
		user := userFromContext(ctx.Request().Context())
		if user == nil {
			return dtos.UserDto{}, fuego.UnauthorizedError{Title: "Unauthorized"}
		}
		return userDto(user), nil
	},
		fuego.OptionOperationID("getCurrentUser"),
	)
}

func NewAuthController(svc *services.AuthService) *AuthController {
	return &AuthController{svc: svc}
}
