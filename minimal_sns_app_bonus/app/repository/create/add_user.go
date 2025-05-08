package create

import (
	"minimal_sns_app/db"
	"minimal_sns_app/interfaces"
	"minimal_sns_app/model"
)

type AddUserHandler struct {
	Repo interfaces.UserRepository
}

func NewAddUserHandler(repo interfaces.UserRepository) *AddUserHandler {
	return &AddUserHandler{Repo: repo}
}

type RealUserRepository struct{}

func (r *RealUserRepository) IsUserIDExists(id string) (bool, error) {
	var count int64
	err := db.DB.Model(&model.User{}).Where("user_id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RealUserRepository) CreateUser(id, name string) error {
	user := model.User{UserID: id, Name: name}
	return db.DB.Create(&user).Error
}
