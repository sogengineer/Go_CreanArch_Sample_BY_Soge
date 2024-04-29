package entity

import "time"

// User is user models property
type User struct {
	UserId    string `gorm:"primaryKey" json:"userId" `
	UserName  string `gorm:"not null" json:"userName,omitempty"`
	Password  string `gorm:"not null" json:"password,omitempty"`
	Email     string `gorm:"not null" json:"email,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
