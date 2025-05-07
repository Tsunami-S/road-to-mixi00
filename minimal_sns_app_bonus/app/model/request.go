package model

type FriendRequest struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	User1ID string `gorm:"type:varchar(20);not null" json:"user1_id"`
	User2ID string `gorm:"type:varchar(20);not null" json:"user2_id"`
	Status  string `gorm:"type:enum('pending','accepted','rejected');default:'pending'" json:"status"`
}

func (FriendRequest) TableName() string {
	return "friend_requests"
}

type AddUserRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type BlockRequest struct {
	User1ID string `json:"user1_id"`
	User2ID string `json:"user2_id"`
}

type FriendRequestInput struct {
	User1ID string `json:"user1_id"`
	User2ID string `json:"user2_id"`
}

type RespondRequestInput struct {
	User1ID string `json:"user1_id"`
	User2ID string `json:"user2_id"`
	Action  string `json:"action"`
}
