package models

import "time"

type Session struct {
	TokenHash string    `gorm:"primaryKey" json:"-"`
	UserID    uint      `gorm:"index;not null" json:"userId"`
	ExpiresAt time.Time `gorm:"not null" json:"expiresAt"`
}

func (Session) TableName() string {
	return "session"
}
