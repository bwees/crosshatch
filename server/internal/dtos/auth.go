package dtos

type LoginDto struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateUserDto struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	IsAdmin  bool   `json:"isAdmin"`
}

type UserDto struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
}

type SetupStatusDto struct {
	SetupRequired bool `json:"setupRequired"`
}
