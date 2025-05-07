package get_all

import (
	"minimal_sns_app/db"
	"minimal_sns_app/model"
)

func Users() ([]model.User, error) {
	var users []model.User
	err := db.DB.Find(&users).Error
	return users, err
}
