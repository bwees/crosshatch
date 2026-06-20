package models

type Printer struct {
	Serial     string `gorm:"primaryKey" json:"serial"`
	Name       string `gorm:"not null" json:"name"`
	HostIP     string `gorm:"not null" json:"host_ip"`
	AccessCode string `gorm:"not null" json:"access_code"`
}
