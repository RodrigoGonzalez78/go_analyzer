package models

import "time"

type Action struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserName    string    `gorm:"not null;index" json:"user_name"`
	Description string    `gorm:"not null" json:"description"`
	Date        time.Time `gorm:"not null" json:"date"`
}
