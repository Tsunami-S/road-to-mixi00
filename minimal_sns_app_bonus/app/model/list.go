package model

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

func (BlockList) TableName() string {
	return "block_list"
}
