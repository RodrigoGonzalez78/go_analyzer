package models

type User struct {
	UserName string `json:"user_name,omitempty" gorm:"primaryKey"`
	Password string `json:"password,omitempty" gorm:"not null"`
}
