package models

type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Username     string `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string `gorm:"not null" json:"-"`
	IsAdmin      bool   `gorm:"not null;default:false" json:"isAdmin"`
}

func (User) TableName() string {
	return "user"
}
