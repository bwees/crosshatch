package dtos

type VapidDto struct {
	PublicKey string `json:"publicKey" validate:"required"`
}

type PushSubscriptionDto struct {
	DeviceId string `json:"deviceId" validate:"required"`
	Endpoint string `json:"endpoint" validate:"required"`
	P256dh   string `json:"p256dh" validate:"required"`
	Auth     string `json:"auth" validate:"required"`
}

type UnsubscribeDto struct {
	DeviceId string `json:"deviceId" validate:"required"`
}

type TestNotificationDto struct {
	DeviceId string `json:"deviceId" validate:"required"`
	Serial   string `json:"serial" validate:"required"`
}

type NotificationSettingsDto struct {
	Enabled        bool `json:"enabled"`
	NotifyComplete bool `json:"notifyComplete"`
	NotifyError    bool `json:"notifyError"`
}
