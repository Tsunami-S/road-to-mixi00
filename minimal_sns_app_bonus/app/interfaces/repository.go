package interfaces

import (
	"errors"
	"gorm.io/gorm"

	"minimal_sns_app/db"
	"minimal_sns_app/model"
)

type FriendRepository interface {
	GetFriends(id string) ([]model.Friend, error)
}

type FriendOfFriendRepository interface {
	GetFriendOfFriend(id string) ([]model.Friend, error)
}

type FriendOfFriendPagingRepository interface {
	GetFriendOfFriendByIDWithPaging(id string, limit, offset int) ([]model.Friend, error)
}

type BlockRepository interface {
	IsBlocked(user1ID, user2ID string) (bool, error)
	DeleteFriendLink(user1ID, user2ID string) error
	RejectRequests(user1ID, user2ID string) error
	Block(user1ID, user2ID string) error
}

type UserRepository interface {
	IsUserIDExists(id string) (bool, error)
	CreateUser(id, name string) error
}

type FriendRequestRepository interface {
	IsBlockedEachOther(user1ID, user2ID string) (bool, error)
	IsAlreadyFriends(user1ID, user2ID string) (bool, error)
	HasPendingRequest(fromID, toID string) (bool, error)
	HasAlreadyRequested(fromID, toID string) (bool, error)
	Request(fromID, toID string) error
	SendRequest(fromID, toID string) error
	RespondRequest(fromID, toID, action string) error
	GetPendingRequests(userID string) ([]model.FriendRequest, error)
}

type FriendRespondRepository interface {
	FindRequest(user1ID, user2ID string) (*model.FriendRequest, error)
	UpdateRequest(req *model.FriendRequest, action string) error
	CreateFriendLink(user1ID, user2ID string) error
}

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
