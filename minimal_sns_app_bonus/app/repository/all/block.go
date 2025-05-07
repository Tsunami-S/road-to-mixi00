package all

import (
	"minimal_sns_app/db"
	"minimal_sns_app/model"
)

func BlockList() ([]model.BlockList, error) {
	var blocks []model.BlockList
	err := db.DB.Find(&blocks).Error
	return blocks, err
}
