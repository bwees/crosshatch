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

type SetPrintSpeedDto struct {
	Level int `json:"level" validate:"required"`
}

type SetFanDto struct {
	Fan   string `json:"fan" validate:"required"`
	Speed int    `json:"speed"`
}

type StartDryingDto struct {
	Temperature int    `json:"temperature" validate:"required"`
	Duration    int    `json:"duration" validate:"required"`
	CoolingTemp int    `json:"coolingTemp"`
	Filament    string `json:"filament"`
	RotateTray  bool   `json:"rotateTray"`
}

type SetFilamentDto struct {
	AmsID         int    `json:"amsId"`
	TrayID        int    `json:"trayId"`
	TrayInfoIdx   string `json:"trayInfoIdx" validate:"required"`
	TrayColor     string `json:"trayColor" validate:"required"`
	TrayType      string `json:"trayType" validate:"required"`
	NozzleTempMin int    `json:"nozzleTempMin" validate:"required"`
	NozzleTempMax int    `json:"nozzleTempMax" validate:"required"`
}
