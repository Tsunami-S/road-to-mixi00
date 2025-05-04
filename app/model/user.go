package model

type User struct {
	ID     uint   `gorm:"primaryKey"`
	UserID string `gorm:"type:varchar(20);unique;not null"`
	Name   string `gorm:"type:varchar(64);not null;default:''"`
}

type Friend struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type FriendLink struct {
	ID      uint   `gorm:"primaryKey"`
	User1ID string `gorm:"type:varchar(20);not null" json:"user1_id"`
	User2ID string `gorm:"type:varchar(20);not null" json:"user2_id"`
}

func (FriendLink) TableName() string {
	return "friend_link"
}

type BlockList struct {
	ID      uint   `gorm:"primaryKey"`
	User1ID string `gorm:"not null" json:"user1_id"`
	User2ID string `gorm:"not null" json:"user2_id"`
}

type FriendRequest struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	User1ID string `gorm:"type:varchar(20);not null" json:"user1_id"`
	User2ID string `gorm:"type:varchar(20);not null" json:"user2_id"`
	Status  string `gorm:"type:enum('pending','accepted','rejected');default:'pending'" json:"status"`
}

func (FriendRequest) TableName() string {
	return "friend_requests"
}
