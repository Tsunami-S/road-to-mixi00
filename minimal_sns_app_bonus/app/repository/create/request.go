package create

import (
	"errors"
	"gorm.io/gorm"

	"minimal_sns_app/db"
	"minimal_sns_app/model"
)

type RealFriendRequestRepository struct{}

func (r *RealFriendRequestRepository) IsBlockedEachOther(user1, user2 string) (bool, error) {
	var block model.BlockList
	err := db.DB.Where(
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		user1, user2, user2, user1,
	).First(&block).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

func (r *RealFriendRequestRepository) IsAlreadyFriends(user1, user2 string) (bool, error) {
	var friend model.FriendLink
	err := db.DB.Where(
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		user1, user2, user2, user1,
	).First(&friend).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

func (r *RealFriendRequestRepository) HasPendingRequest(user1, user2 string) (bool, error) {
	var req model.FriendRequest
	err := db.DB.Where("user1_id = ? AND user2_id = ? AND status = 'pending'", user1, user2).First(&req).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

func (r *RealFriendRequestRepository) Request(user1, user2 string) error {
	req := model.FriendRequest{
		User1ID: user1,
		User2ID: user2,
		Status:  "pending",
	}
	return db.DB.Create(&req).Error
}

func (r *RealFriendRequestRepository) GetPendingRequests(userID string) ([]model.FriendRequest, error) {
	var requests []model.FriendRequest
	err := db.DB.
		Where("user2_id = ? AND status = 'pending'", userID).
		Find(&requests).Error
	return requests, err
}

func (r *RealFriendRequestRepository) HasAlreadyRequested(fromID, toID string) (bool, error) {
	var req model.FriendRequest
	err := db.DB.Where("user1_id = ? AND user2_id = ? AND status = ?", fromID, toID, "pending").
		First(&req).Error

	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}
