package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UserName string `gorm:"unique;not null" json:"user_name"`
	Password string `gorm:"not null" json:"password"`
}
