package utils

import (
	"zadanie-6105/src/main/models"

	"gorm.io/gorm"
)

// Проверяем существует ли пользователь (проверка происходит по имени)
func UserExists(db *gorm.DB, username string) (bool, error) {
	var count int64
	err := db.Table("employee").Model(&models.Employee{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Проверяем существует ли пользователь (проверка происходит по id)
func UserIdExists(db *gorm.DB, id string) (bool, error) {
	var count int64
	err := db.Table("employee").Model(&models.Employee{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Проверяем существует ли тендер (по id)
func TenderExists(db *gorm.DB, id string) (bool, error) {
	var count int64
	err := db.Table("tender").Model(&models.Tender{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
