package dtos

type CreatePrinterDto struct {
	Serial     string `json:"serial" validate:"required"`
	Name       string `json:"name" validate:"required"`
	HostIP     string `json:"hostIp" validate:"required"`
	AccessCode string `json:"accessCode" validate:"required"`
}

type UpdatePrinterDto struct {
	Name       *string `json:"name,omitempty"`
	HostIP     *string `json:"hostIp,omitempty"`
	AccessCode *string `json:"accessCode,omitempty"`
}

type SetLightDto struct {
	State bool `json:"state"`
}
