package create

import (
	"testing"

	"minimal_sns_app/db"
	"minimal_sns_app/model"

	"github.com/stretchr/testify/assert"
)

func setupTestDB_Block(t *testing.T) {
	db.InitDB()
}

func TestIsBlocked(t *testing.T) {
	setupTestDB_Block(t)
	repo := &RealBlockRepository{}

	t.Run("ブロックあり", func(t *testing.T) {
		ok, err := repo.IsBlocked("id3", "id2")
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("ブロックなし", func(t *testing.T) {
		ok, err := repo.IsBlocked("id1", "id2")
		assert.NoError(t, err)
		assert.False(t, ok)
	})
}

func TestBlock(t *testing.T) {
	setupTestDB_Block(t)
	repo := &RealBlockRepository{}

	err := repo.Block("id1", "id2")
	assert.NoError(t, err)

	var blk model.BlockList
	err = db.DB.Where("user1_id = ? AND user2_id = ?", "id1", "id2").First(&blk).Error
	assert.NoError(t, err)
}

func TestDeleteFriendLink(t *testing.T) {
	setupTestDB_Block(t)
	repo := &RealBlockRepository{}

	err := repo.DeleteFriendLink("id1", "id2")
	assert.NoError(t, err)

	var link model.FriendLink
	err = db.DB.Where("user1_id = ? AND user2_id = ?", "id1", "id2").First(&link).Error
	assert.Error(t, err)
}

func TestRejectRequests(t *testing.T) {
	setupTestDB_Block(t)
	repo := &RealBlockRepository{}

	err := repo.RejectRequests("id1", "id2")
	assert.NoError(t, err)

	var req model.FriendRequest
	db.DB.Where("user1_id = ? AND user2_id = ?", "id2", "id1").First(&req)
	assert.Equal(t, "rejected", req.Status)
}
