package get

import (
	"testing"

	"minimal_sns_app/db"

	"github.com/stretchr/testify/assert"
)

func setupTestDB_PendingRequest(t *testing.T) {
	db.InitDB()
}

func TestGetPendingRequests(t *testing.T) {
	setupTestDB_PendingRequest(t)

	repo := &RealFriendRequestRepository{}

	t.Run("正常系: id1 への保留中リクエスト", func(t *testing.T) {
		requests, err := repo.GetPendingRequests("id27")
		assert.NoError(t, err)

		var gotIDs []string
		for _, r := range requests {
			gotIDs = append(gotIDs, r.User1ID)
		}

		assert.Contains(t, gotIDs, "id1")
		assert.NotContains(t, gotIDs, "id3")
	})
}
