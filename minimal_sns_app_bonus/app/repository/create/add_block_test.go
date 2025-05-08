package create

import (
	"testing"

	"minimal_sns_app/db"
	"minimal_sns_app/model"

	"github.com/stretchr/testify/assert"
)

func setupBlockTestDB(t *testing.T) {
	db.InitDB()
	db.DB.Exec("DELETE FROM block_lists")
	db.DB.Exec("DELETE FROM friend_links")
	db.DB.Exec("DELETE FROM friend_requests")
	db.DB.Exec("DELETE FROM users")

	// ダミーユーザー作成
	db.DB.Create(&model.User{UserID: "id1", Name: "User 1"})
	db.DB.Create(&model.User{UserID: "id2", Name: "User 2"})
	db.DB.Create(&model.User{UserID: "id3", Name: "User 3"})

	// ブロック・フレンド・リクエストデータ
	db.DB.Create(&model.BlockList{User1ID: "id3", User2ID: "id2"})
	db.DB.Create(&model.FriendLink{User1ID: "id1", User2ID: "id2"})
	db.DB.Create(&model.FriendRequest{User1ID: "id2", User2ID: "id1", Status: "pending"})
}

func TestRealBlockRepository(t *testing.T) {
	setupBlockTestDB(t)
	repo := &RealBlockRepository{}

	t.Run("IsBlocked - true", func(t *testing.T) {
		ok, err := repo.IsBlocked("id3", "id2")
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("IsBlocked - false", func(t *testing.T) {
		ok, err := repo.IsBlocked("id1", "id2")
		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("Block - success", func(t *testing.T) {
		err := repo.Block("id1", "id3")
		assert.NoError(t, err)

		var block model.BlockList
		err = db.DB.Where("user1_id = ? AND user2_id = ?", "id1", "id3").First(&block).Error
		assert.NoError(t, err)
	})

	t.Run("DeleteFriendLink - success", func(t *testing.T) {
		err := repo.DeleteFriendLink("id1", "id2")
		assert.NoError(t, err)

		var link model.FriendLink
		err = db.DB.Where("user1_id = ? AND user2_id = ?", "id1", "id2").First(&link).Error
		assert.Error(t, err) // 削除済みなのでエラーになる
	})

	t.Run("RejectRequests - success", func(t *testing.T) {
		err := repo.RejectRequests("id1", "id2")
		assert.NoError(t, err)

		var req model.FriendRequest
		err = db.DB.Where("user1_id = ? AND user2_id = ?", "id2", "id1").First(&req).Error
		assert.NoError(t, err)
		assert.Equal(t, "rejected", req.Status)
	})
}
