package models

type User struct {
	UserName string `gorm:"primaryKey" json:"user_name"`
	Password string `gorm:"not null" json:"password"`
}
