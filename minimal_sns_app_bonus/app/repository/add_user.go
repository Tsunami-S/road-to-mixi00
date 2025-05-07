package repository

import (
	"minimal_sns_app/db"
	"minimal_sns_app/model"
)

func IsUserIDExists(id string) (bool, error) {
	var count int64
	err := db.DB.Model(&model.User{}).Where("user_id = ?", id).Count(&count).Error
	return count > 0, err
}

func CreateUser(id, name string) error {
	user := model.User{UserID: id, Name: name}
	return db.DB.Create(&user).Error
}
