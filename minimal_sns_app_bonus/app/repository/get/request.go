package get

import (
	"errors"

	"gorm.io/gorm"
	"minimal_sns_app/db"
	"minimal_sns_app/model"
)

type RealFriendRequestRepository struct{}

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

func (r *RealFriendRequestRepository) HasAlreadyRequested(user1, user2 string) (bool, error) {
	var req model.FriendRequest
	err := db.DB.
		Where("user1_id = ? AND user2_id = ? AND status = ?", user1, user2, "pending").
		First(&req).Error

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

func (r *RealFriendRequestRepository) GetPendingRequests(userID string) ([]model.FriendRequest, error) {
	var requests []model.FriendRequest

	query := `
	SELECT fr.*
	FROM friend_requests fr
	WHERE fr.user2_id = ?
	  AND fr.status = 'pending'
	  AND fr.user1_id != fr.user2_id
	  AND NOT EXISTS (
	    SELECT 1 FROM block_list b
	    WHERE 
	      (b.user1_id = fr.user1_id AND b.user2_id = fr.user2_id)
	      OR (b.user1_id = fr.user2_id AND b.user2_id = fr.user1_id)
	  );
	`

	err := db.DB.Raw(query, userID).Scan(&requests).Error
	return requests, err
}

func (r *RealFriendRequestRepository) Request(fromID, toID string) error {
	return nil
}

func (r *RealFriendRequestRepository) RespondRequest(fromID, toID, action string) error {
	return nil
}

func (r *RealFriendRequestRepository) SendRequest(fromID, toID string) error {
	return nil
}

func (r *RealFriendRequestRepository) IsBlockedEachOther(user1ID, user2ID string) (bool, error) {
	return false, nil
}
