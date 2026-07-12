package controllers

import (
	"strconv"

	"crosshatch/internal/dtos"
	"crosshatch/internal/services"

	"github.com/go-fuego/fuego"
)

type UsersController struct {
	svc *services.AuthService
}

func (c *UsersController) requireAdmin(ctx fuego.ContextNoBody) error {
	user := userFromContext(ctx.Request().Context())
	if user == nil || !user.IsAdmin {
		return services.ErrForbidden
	}
	return nil
}

func (c *UsersController) Register(api *fuego.Server) {
	route := fuego.Group(api, "/users")

	fuego.Get(route, "/", func(ctx fuego.ContextNoBody) ([]dtos.UserDto, error) {
		if err := c.requireAdmin(ctx); err != nil {
			return nil, err
		}

		users, err := c.svc.ListUsers()
		if err != nil {
			return nil, err
		}

		out := make([]dtos.UserDto, len(users))
		for i, u := range users {
			out[i] = userDto(&u)
		}
		return out, nil
	},
		fuego.OptionOperationID("getUsers"),
	)

	fuego.Post(route, "/", func(ctx fuego.ContextWithBody[dtos.CreateUserDto]) (dtos.UserDto, error) {
		user := userFromContext(ctx.Request().Context())
		if user == nil || !user.IsAdmin {
			return dtos.UserDto{}, services.ErrForbidden
		}

		dto, err := ctx.Body()
		if err != nil {
			return dtos.UserDto{}, err
		}

		created, err := c.svc.CreateUser(dto)
		if err != nil {
			return dtos.UserDto{}, err
		}
		return userDto(created), nil
	},
		fuego.OptionOperationID("createUser"),
		fuego.OptionDefaultStatusCode(201),
	)

	fuego.Delete(route, "/{id}", func(ctx fuego.ContextNoBody) (any, error) {
		if err := c.requireAdmin(ctx); err != nil {
			return nil, err
		}

		id, err := strconv.ParseUint(ctx.PathParam("id"), 10, 64)
		if err != nil {
			return nil, fuego.BadRequestError{Title: "Bad Request", Detail: "invalid user id"}
		}

		current := userFromContext(ctx.Request().Context())
		if current != nil && uint(id) == current.ID {
			return nil, fuego.BadRequestError{Title: "Bad Request", Detail: "cannot delete your own account"}
		}

		return nil, c.svc.DeleteUser(uint(id))
	},
		fuego.OptionOperationID("deleteUser"),
		fuego.OptionDefaultStatusCode(204),
	)
}

func NewUsersController(svc *services.AuthService) *UsersController {
	return &UsersController{svc: svc}
}
