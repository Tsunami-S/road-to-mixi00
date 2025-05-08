package create

import (
	"errors"

	"gorm.io/gorm"
	"minimal_sns_app/db"
	"minimal_sns_app/interfaces"
	"minimal_sns_app/model"
)

type RespondHandler struct {
	Validator interfaces.Validator
	Repo      interfaces.FriendRequestRepository
}

func NewRespondRequestHandler(v interfaces.Validator, r interfaces.FriendRequestRepository) *RespondHandler {
	return &RespondHandler{Validator: v, Repo: r}
}

func (r *RealFriendRequestRepository) UpdateRequest(req *model.FriendRequest, status string) error {
	return db.DB.Model(req).Update("status", status).Error
}

func (r *RealFriendRequestRepository) FriendLink(user1, user2 string) error {
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
