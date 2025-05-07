package create

import (
	"minimal_sns_app/db"
	"minimal_sns_app/model"
)

func FindRequest(user1ID, user2ID string) (*model.FriendRequest, error) {
	var req model.FriendRequest
	err := db.DB.Where("user1_id = ? AND user2_id = ? AND status = 'pending'", user1ID, user2ID).First(&req).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func UpdateRequest(req *model.FriendRequest, status string) error {
	return db.DB.Model(req).Update("status", status).Error
}

func FriendLink(user1ID, user2ID string) error {
	link := model.FriendLink{User1ID: user1ID, User2ID: user2ID}
	return db.DB.Create(&link).Error
}
