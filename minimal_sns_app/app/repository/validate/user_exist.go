package validate

import (
	"gorm.io/gorm"
	"minimal_sns_app/db"
	"minimal_sns_app/model"
)

func UserExists(id int) (bool, error) {
	var user model.User
	err := db.DB.Where("user_id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
