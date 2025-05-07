package create

import (
	"minimal_sns_app/db"
	"minimal_sns_app/model"
)

func IsBlocked(user1, user2 string) (bool, error) {
	var block model.BlockList
	err := db.DB.Where("user1_id = ? AND user2_id = ?", user1, user2).First(&block).Error
	if err == nil {
		return true, nil
	}
	if err.Error() == "record not found" {
		return false, nil
	}
	return false, err
}

func Block(user1, user2 string) error {
	block := model.BlockList{User1ID: user1, User2ID: user2}
	return db.DB.Create(&block).Error
}

func DeleteFriendLink(user1, user2 string) error {
	return db.DB.Where(
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		user1, user2, user2, user1,
	).Delete(&model.FriendLink{}).Error
}

func RejectRequests(user1, user2 string) error {
	return db.DB.Model(&model.FriendRequest{}).
		Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)", user1, user2, user2, user1).
		Where("status = ?", "pending").
		Update("status", "rejected").Error
}
