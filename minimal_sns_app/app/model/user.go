package model

type User struct {
	ID     uint   `gorm:"primaryKey"`
	UserID int    `gorm:not null"`
	Name   string `gorm:"type:varchar(64);not null;default:''"`
}

type Friend struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type FriendLink struct {
	ID      uint `gorm:"primaryKey"`
	User1ID int  `gorm:not null" json:"user1_id"`
	User2ID int  `gorm:not null" json:"user2_id"`
}

func (FriendLink) TableName() string {
	return "friend_link"
}

type BlockList struct {
	ID      uint `gorm:"primaryKey"`
	User1ID int  `gorm:"not null" json:"user1_id"`
	User2ID int  `gorm:"not null" json:"user2_id"`
}

func (BlockList) TableName() string {
	return "block_list"
}
