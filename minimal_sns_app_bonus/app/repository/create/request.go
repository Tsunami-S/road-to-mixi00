package create

import (
	"minimal_sns_app/db"
	"minimal_sns_app/model"
)

func IsBlockedEachOther(user1, user2 string) (bool, error) {
	var block model.BlockList
	err := db.DB.Where(
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		user1, user2, user2, user1,
	).First(&block).Error
	if err == nil {
		return true, nil
	}
	if err.Error() == "record not found" {
		return false, nil
	}
	return false, err
}

func IsAlreadyFriends(user1, user2 string) (bool, error) {
	var friend model.FriendLink
	err := db.DB.Where(
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		user1, user2, user2, user1,
	).First(&friend).Error
	if err == nil {
		return true, nil
	}
	if err.Error() == "record not found" {
		return false, nil
	}
	return false, err
}

func HasPendingRequest(user1, user2 string) (bool, error) {
	var req model.FriendRequest
	err := db.DB.Where("user1_id = ? AND user2_id = ? AND status = 'pending'", user1, user2).First(&req).Error
	if err == nil {
		return true, nil
	}
	if err.Error() == "record not found" {
		return false, nil
	}
	return false, err
}

func HasAlreadyRequested(user1, user2 string) (bool, error) {
	var req model.FriendRequest
	err := db.DB.Where("user1_id = ? AND user2_id = ? AND status = 'pending'", user1, user2).First(&req).Error
	if err == nil {
		return true, nil
	}
	if err.Error() == "record not found" {
		return false, nil
	}
	return false, err
}

func Request(user1, user2 string) error {
	req := model.FriendRequest{User1ID: user1, User2ID: user2, Status: "pending"}
	return db.DB.Create(&req).Error
}
