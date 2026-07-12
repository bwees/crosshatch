package models

import (
	"time"

	"crosshatch/internal/dtos"
)

type PushSubscription struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID  string    `gorm:"uniqueIndex;not null" json:"deviceId"`
	UserID    string    `gorm:"index;not null" json:"userId"`
	Endpoint  string    `gorm:"not null" json:"endpoint"`
	P256dh    string    `json:"p256dh"`
	Auth      string    `json:"auth"`
	CreatedAt time.Time `json:"createdAt"`
}

func (PushSubscription) TableName() string {
	return "push_subscription"
}

type NotificationSetting struct {
	DeviceID       string `gorm:"primaryKey" json:"deviceId"`
	PrinterSerial  string `gorm:"primaryKey" json:"printerSerial"`
	UserID         string `gorm:"index;not null" json:"userId"`
	Enabled        bool   `gorm:"not null;default:false" json:"enabled"`
	NotifyComplete bool   `gorm:"not null;default:false" json:"notifyComplete"`
	NotifyError    bool   `gorm:"not null;default:false" json:"notifyError"`
}

func (NotificationSetting) TableName() string {
	return "notification_setting"
}

func (s NotificationSetting) ToDto() dtos.NotificationSettingsDto {
	return dtos.NotificationSettingsDto{
		Enabled:        s.Enabled,
		NotifyComplete: s.NotifyComplete,
		NotifyError:    s.NotifyError,
	}
}

type AppConfig struct {
	Key   string `gorm:"primaryKey" json:"key"`
	Value string `gorm:"not null" json:"value"`
}

func (AppConfig) TableName() string {
	return "app_config"
}
