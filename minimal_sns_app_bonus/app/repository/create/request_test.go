package create

import (
	"testing"

	"minimal_sns_app/db"
	"minimal_sns_app/model"

	"github.com/stretchr/testify/assert"
)

func setupFriendStatusTestDB(t *testing.T) {
	db.InitDB()
}

func TestFriendStatusFunctions(t *testing.T) {
	setupFriendStatusTestDB(t)

	t.Run("IsBlockedEachOther returns true when blocked", func(t *testing.T) {
		ok, err := IsBlockedEachOther("id1", "id39")
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("IsBlockedEachOther returns false when not blocked", func(t *testing.T) {
		ok, err := IsBlockedEachOther("id1", "id5")
		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("IsAlreadyFriends returns true when friends", func(t *testing.T) {
		ok, err := IsAlreadyFriends("id1", "id2")
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("IsAlreadyFriends returns false when not friends", func(t *testing.T) {
		ok, err := IsAlreadyFriends("id1", "id44")
		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("HasPendingRequest returns true when pending exists", func(t *testing.T) {
		ok, err := HasPendingRequest("id1", "id41")
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("HasPendingRequest returns false when no pending", func(t *testing.T) {
		ok, err := HasPendingRequest("id1", "id44")
		assert.NoError(t, err)
		assert.False(t, ok)
	})
}

func TestHasAlreadyRequested(t *testing.T) {
	setupFriendStatusTestDB(t)

	t.Run("正常系: 既にリクエストあり", func(t *testing.T) {
		_ = db.DB.Create(&model.FriendRequest{User1ID: "id1", User2ID: "id2", Status: "pending"}).Error

		found, err := HasAlreadyRequested("id1", "id2")
		assert.NoError(t, err)
		assert.True(t, found)
	})

	t.Run("異常系: リクエストが存在しない", func(t *testing.T) {
		found, err := HasAlreadyRequested("no_user1", "no_user2")
		assert.NoError(t, err)
		assert.False(t, found)
	})
}

func TestRequest(t *testing.T) {
	setupFriendStatusTestDB(t)

	t.Run("正常系: 新しいリクエストを作成", func(t *testing.T) {
		err := Request("id3", "id4")
		assert.NoError(t, err)

		var req model.FriendRequest
		err = db.DB.Where("user1_id = ? AND user2_id = ?", "id3", "id4").First(&req).Error
		assert.NoError(t, err)
		assert.Equal(t, "pending", req.Status)
	})
}
