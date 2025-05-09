package create

import (
	"testing"

	"minimal_sns_app/db"
	"minimal_sns_app/model"

	"github.com/stretchr/testify/assert"
)

func setupTestDB_FriendRequestRepo(t *testing.T) {
	db.InitDB()
}

func TestRealFriendRequestRepository_IsBlockedEachOther(t *testing.T) {
	setupTestDB_FriendRequestRepo(t)
	repo := &RealFriendRequestRepository{}

	blocked, err := repo.IsBlockedEachOther("id1", "id39")
	assert.NoError(t, err)
	assert.True(t, blocked)

	notBlocked, err := repo.IsBlockedEachOther("id46", "id49")
	assert.NoError(t, err)
	assert.False(t, notBlocked)
}

func TestRealFriendRequestRepository_IsAlreadyFriends(t *testing.T) {
	setupTestDB_FriendRequestRepo(t)
	repo := &RealFriendRequestRepository{}

	isFriend, err := repo.IsAlreadyFriends("id1", "id2")
	assert.NoError(t, err)
	assert.True(t, isFriend)

	notFriend, err := repo.IsAlreadyFriends("id1", "id49")
	assert.NoError(t, err)
	assert.False(t, notFriend)
}

func TestRealFriendRequestRepository_HasPendingRequest(t *testing.T) {
	setupTestDB_FriendRequestRepo(t)
	repo := &RealFriendRequestRepository{}

	hasPending, err := repo.HasPendingRequest("id1", "id41")
	assert.NoError(t, err)
	assert.True(t, hasPending)

	noPending, err := repo.HasPendingRequest("id2", "id4")
	assert.NoError(t, err)
	assert.False(t, noPending)
}

func TestRealFriendRequestRepository_Request(t *testing.T) {
	setupTestDB_FriendRequestRepo(t)
	repo := &RealFriendRequestRepository{}

	err := repo.Request("id10", "id13")
	assert.NoError(t, err)

	var reqExist bool
	err = db.DB.Model(&model.FriendRequest{}).
		Select("count(*) > 0").
		Where("user1_id = ? AND user2_id = ? AND status = ?", "id10", "id13", "pending").
		Find(&reqExist).Error
	assert.NoError(t, err)
	assert.True(t, reqExist)
}

func TestRealFriendRequestRepository_HasAlreadyRequested(t *testing.T) {
	setupTestDB_FriendRequestRepo(t)
	repo := &RealFriendRequestRepository{}

	hasRequested, err := repo.HasAlreadyRequested("id1", "id41")
	assert.NoError(t, err)
	assert.True(t, hasRequested)

	notRequested, err := repo.HasAlreadyRequested("id2", "id3")
	assert.NoError(t, err)
	assert.False(t, notRequested)
}

func TestRealFriendRequestRepository_GetPendingRequests(t *testing.T) {
	setupTestDB_FriendRequestRepo(t)
	repo := &RealFriendRequestRepository{}

	reqs, err := repo.GetPendingRequests("id27")
	assert.NoError(t, err)
	assert.Greater(t, len(reqs), 0)

	reqs, err = repo.GetPendingRequests("id2")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(reqs))
}
