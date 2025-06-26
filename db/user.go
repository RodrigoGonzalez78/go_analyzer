package db

import (
	"fmt"

	"github.com/RodrigoGonzalez78/go_analyzer/models"
	"gorm.io/gorm"
)

func CreateUser(user models.User) error {
	// La verificación de nombre de usuario único ya se hace en el endpoint
	if err := database.Create(&user).Error; err != nil {
		return fmt.Errorf("error al crear usuario: %v", err)
	}

	return nil
}

func IsUserNameUnique(userName string) (bool, error) {
	var count int64
	err := database.Model(&models.User{}).Where("user_name = ?", userName).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func GetUserByUserName(userName string) (*models.User, error) {
	var user models.User
	err := database.Model(&models.User{}).Where("user_name = ?", userName).First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
