package get_all

import (
	"minimal_sns_app/db"
	"minimal_sns_app/model"
)

func FriendRequests() ([]model.FriendRequest, error) {
	var requests []model.FriendRequest
	err := db.DB.Find(&requests).Error
	return requests, err
}
