package repository

import (
	"minimal_sns_app/db"
	"minimal_sns_app/model"
)

func GetAllFriendRequests() ([]model.FriendRequest, error) {
	var requests []model.FriendRequest
	err := db.DB.Find(&requests).Error
	return requests, err
}
