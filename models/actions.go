package models

type Action struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	UserID  uint   `gorm:"not null;index" json:"user_id"`
	Content string `gorm:"not null" json:"content"`
}
