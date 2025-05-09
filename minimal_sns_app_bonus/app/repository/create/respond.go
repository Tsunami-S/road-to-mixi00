package create

import (
	"errors"
	"gorm.io/gorm"

	"minimal_sns_app/db"
	"minimal_sns_app/model"
)

type RealFriendRespondRepository struct{}

func (r *RealFriendRespondRepository) RespondRequest(fromID, toID, action string) error {
	return db.DB.Model(&model.FriendRequest{}).
		Where("user1_id = ? AND user2_id = ?", fromID, toID).
		Update("status", action).Error
}

func (r *RealFriendRespondRepository) FindRequest(user1, user2 string) (*model.FriendRequest, error) {
	var req model.FriendRequest
	err := db.DB.Where("user1_id = ? AND user2_id = ? AND status = 'pending'", user1, user2).First(&req).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("request not found or already handled")
	}
	return &req, err
}

func (r *RealFriendRespondRepository) UpdateRequest(req *model.FriendRequest, status string) error {
	return db.DB.Model(req).Update("status", status).Error
}

func (r *RealFriendRespondRepository) CreateFriendLink(user1, user2 string) error {
	var existing model.FriendLink
	err := db.DB.Where(
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		user1, user2, user2, user1,
	).First(&existing).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	link := model.FriendLink{User1ID: user1, User2ID: user2}
	return db.DB.Create(&link).Error
}
