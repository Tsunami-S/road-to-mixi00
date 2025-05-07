package all

import (
	"minimal_sns_app/db"
	"minimal_sns_app/model"
)

func GetAllFriendLinks() ([]model.FriendLink, error) {
	var links []model.FriendLink
	err := db.DB.Find(&links).Error
	return links, err
}
