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
