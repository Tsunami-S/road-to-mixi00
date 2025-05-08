package create

import (
	"testing"

	"minimal_sns_app/db"
	"minimal_sns_app/model"

	"github.com/stretchr/testify/assert"
)

func setupTestDB_Respond(t *testing.T) {
	db.InitDB()
}

func TestFindRequest(t *testing.T) {
	setupTestDB_Respond(t)
	repo := &RealFriendRespondRepository{}

	t.Run("正常系: リクエスト存在", func(t *testing.T) {
		req, err := repo.FindRequest("id1", "id41")
		assert.NoError(t, err)
		assert.NotNil(t, req)
	})

	t.Run("異常系: リクエストなし", func(t *testing.T) {
		_, err := repo.FindRequest("id3", "id9")
		assert.Error(t, err)
	})
}

func TestUpdateRequest(t *testing.T) {
	setupTestDB_Respond(t)
	repo := &RealFriendRespondRepository{}

	req, err := repo.FindRequest("id1", "id41")
	assert.NoError(t, err)
	assert.NotNil(t, req)

	err = repo.UpdateRequest(req, "accepted")
	assert.NoError(t, err)

	db.DB.First(&req, req.ID)
	assert.Equal(t, "accepted", req.Status)
}

func TestCreateFriendLink(t *testing.T) {
	setupTestDB_Respond(t)
	repo := &RealFriendRespondRepository{}

	err := repo.CreateFriendLink("id1", "id2")
	assert.NoError(t, err)

	var link model.FriendLink
	err = db.DB.Where("user1_id = ? AND user2_id = ?", "id1", "id2").First(&link).Error
	assert.NoError(t, err)
}
