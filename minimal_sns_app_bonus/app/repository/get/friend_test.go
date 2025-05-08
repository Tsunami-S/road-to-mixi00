package get

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"minimal_sns_app/db"
)

func setupTestDB_FriendRepo(t *testing.T) {
	db.InitDB()
}

func TestFriendRepository_Friend(t *testing.T) {
	setupTestDB_FriendRepo(t)
	repo := &RealFriendRepository{}

	tests := []struct {
		name       string
		userID     string
		shouldHave []string
		shouldNot  []string
	}{
		{
			name:       "正常系: id32 の友達取得",
			userID:     "id32",
			shouldHave: []string{"id5"},
			shouldNot:  []string{"id7"},
		},
		{
			name:       "存在しないユーザー",
			userID:     "invalid_id",
			shouldHave: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			friends, err := repo.GetFriends(tc.userID)
			assert.NoError(t, err)

			gotIDs := map[string]bool{}
			for _, f := range friends {
				gotIDs[f.ID] = true
			}

			for _, id := range tc.shouldHave {
				assert.True(t, gotIDs[id], id+" は含まれるべき")
			}
			for _, id := range tc.shouldNot {
				assert.False(t, gotIDs[id], id+" は除外されるべき")
			}
		})
	}
}
