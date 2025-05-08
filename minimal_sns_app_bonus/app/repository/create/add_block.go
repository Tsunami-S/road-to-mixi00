package create

import (
	"errors"

	"gorm.io/gorm"
	"minimal_sns_app/db"
	"minimal_sns_app/model"
)

type RealBlockRepository struct{}

func (r *RealBlockRepository) IsBlocked(user1, user2 string) (bool, error) {
	var block model.BlockList
	err := db.DB.Where("user1_id = ? AND user2_id = ?", user1, user2).First(&block).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *RealBlockRepository) Block(user1, user2 string) error {
	blocked, err := r.IsBlocked(user1, user2)
	if err != nil {
		return err
	}
	if blocked {
		return nil
	}
	block := model.BlockList{User1ID: user1, User2ID: user2}
	return db.DB.Create(&block).Error
}

func (r *RealBlockRepository) DeleteFriendLink(user1, user2 string) error {
	return db.DB.Where(
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		user1, user2, user2, user1,
	).Delete(&model.FriendLink{}).Error
}

func (r *RealBlockRepository) RejectRequests(user1, user2 string) error {
	return db.DB.Model(&model.FriendRequest{}).
		Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)", user1, user2, user2, user1).
		Where("status = ?", "pending").
		Update("status", "rejected").Error
}
