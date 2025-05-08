package create

import (
	"testing"

	"minimal_sns_app/db"
	"minimal_sns_app/model"

	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) {
	db.InitDB()
}

func TestFindRequest(t *testing.T) {
	setupTestDB(t)

	db.DB.Create(&model.FriendRequest{User1ID: "id1", User2ID: "id2", Status: "pending"})

	t.Run("pending request found", func(t *testing.T) {
		req, err := FindRequest("id1", "id2")
		assert.NoError(t, err)
		assert.NotNil(t, req)
		assert.Equal(t, "id1", req.User1ID)
	})

	t.Run("request not found", func(t *testing.T) {
		req, err := FindRequest("nonexistent", "user")
		assert.NoError(t, err)
		assert.Nil(t, req)
	})
}

func TestUpdateRequest(t *testing.T) {
	setupTestDB(t)

	req := &model.FriendRequest{User1ID: "id3", User2ID: "id4", Status: "pending"}
	db.DB.Create(req)

	err := UpdateRequest(req, "accepted")
	assert.NoError(t, err)

	var updated model.FriendRequest
	db.DB.First(&updated, req.ID)
	assert.Equal(t, "accepted", updated.Status)
}

func TestFriendLink(t *testing.T) {
	setupTestDB(t)

	db.DB.Where("user1_id = ? AND user2_id = ?", "id5", "id6").Delete(&model.FriendLink{})

	t.Run("create new link", func(t *testing.T) {
		err := FriendLink("id5", "id6")
		assert.NoError(t, err)

		var link model.FriendLink
		result := db.DB.Where("user1_id = ? AND user2_id = ?", "id5", "id6").First(&link)
		assert.NoError(t, result.Error)
	})

	t.Run("link already exists", func(t *testing.T) {
		err := FriendLink("id5", "id6")
		assert.NoError(t, err)
	})
}
