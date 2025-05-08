package create

import (
	"testing"

	"minimal_sns_app/db"
	"minimal_sns_app/model"

	"github.com/stretchr/testify/assert"
)

func setupTestDB_FriendRequestRepo(t *testing.T) {
	db.InitDB()
	db.DB.Exec("DELETE FROM friend_requests")
	db.DB.Exec("DELETE FROM friend_links")
	db.DB.Exec("DELETE FROM block_lists")
	db.DB.Exec("DELETE FROM users")

	db.DB.Create(&model.User{UserID: "id1", Name: "User 1"})
	db.DB.Create(&model.User{UserID: "id2", Name: "User 2"})
	db.DB.Create(&model.User{UserID: "id3", Name: "User 3"})
	db.DB.Create(&model.User{UserID: "id4", Name: "User 4"})

	db.DB.Create(&model.FriendLink{User1ID: "id1", User2ID: "id2"})
	db.DB.Create(&model.BlockList{User1ID: "id3", User2ID: "id1"})
	db.DB.Create(&model.FriendRequest{User1ID: "id1", User2ID: "id4", Status: "pending"})
}

func TestRealFriendRequestRepository_IsBlockedEachOther(t *testing.T) {
	setupTestDB_FriendRequestRepo(t)
	repo := &RealFriendRequestRepository{}

	blocked, err := repo.IsBlockedEachOther("id20", "id1")
	assert.NoError(t, err)
	assert.True(t, blocked)

	notBlocked, err := repo.IsBlockedEachOther("id1", "id2")
	assert.NoError(t, err)
	assert.False(t, notBlocked)
}

func TestRealFriendRequestRepository_IsAlreadyFriends(t *testing.T) {
	setupTestDB_FriendRequestRepo(t)
	repo := &RealFriendRequestRepository{}

	isFriend, err := repo.IsAlreadyFriends("id1", "id2")
	assert.NoError(t, err)
	assert.True(t, isFriend)

	notFriend, err := repo.IsAlreadyFriends("id1", "id40")
	assert.NoError(t, err)
	assert.False(t, notFriend)
}

func TestRealFriendRequestRepository_HasPendingRequest(t *testing.T) {
	setupTestDB_FriendRequestRepo(t)
	repo := &RealFriendRequestRepository{}

	hasPending, err := repo.HasPendingRequest("id1", "id4")
	assert.NoError(t, err)
	assert.True(t, hasPending)

	noPending, err := repo.HasPendingRequest("id2", "id4")
	assert.NoError(t, err)
	assert.False(t, noPending)
}

func TestRealFriendRequestRepository_Request(t *testing.T) {
	setupTestDB_FriendRequestRepo(t)
	repo := &RealFriendRequestRepository{}

	err := repo.Request("id2", "id3")
	assert.NoError(t, err)

	var req model.FriendRequest
	err = db.DB.Where("user1_id = ? AND user2_id = ?", "id2", "id3").First(&req).Error
	assert.NoError(t, err)
	assert.Equal(t, "pending", req.Status)
}
