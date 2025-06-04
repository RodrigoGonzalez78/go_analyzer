package db

import (
	"fmt"

	"github.com/RodrigoGonzalez78/go_analyzer/models"
)

func CreateAction(action models.Action) error {
	if err := database.Create(&action).Error; err != nil {
		return fmt.Errorf("error al crear la accion: %v", err)
	}
	return nil
}

func GetActionByID(id uint) (*models.Action, error) {
	var action models.Action
	err := database.First(&action, id).Error
	if err != nil {
		return nil, err
	}
	return &action, nil
}

func GetUserActionsPaginated(userName string, page int, pageSize int) ([]models.Action, error) {
	var actions []models.Action
	offset := (page - 1) * pageSize

	err := database.
		Where("user_name = ?", userName).
		Limit(pageSize).
		Offset(offset).
		Find(&actions).Error

	if err != nil {
		return nil, err
	}
	return actions, nil
}

func DeleteActionByID(id uint) error {
	result := database.Delete(&models.Action{}, id)
	return result.Error
}
