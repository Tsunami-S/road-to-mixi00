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
