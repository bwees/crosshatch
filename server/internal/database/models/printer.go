package models

import "fmt"

type Printer struct {
	Serial     string `gorm:"primaryKey" json:"serial" validate:"required"`
	Name       string `gorm:"not null" json:"name" validate:"required"`
	HostIP     string `gorm:"not null" json:"hostIp" validate:"required"`
	AccessCode string `gorm:"not null" json:"accessCode" validate:"required"`
}

func (p *Printer) CameraURL() string {
	return fmt.Sprintf("rtsps://bblp:%s@%s:322/streaming/live/1", p.AccessCode, p.HostIP)
}

func (Printer) TableName() string {
	return "printer"
}
