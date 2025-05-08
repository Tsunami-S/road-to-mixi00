package get

import (
	"testing"

	"minimal_sns_app/db"

	"github.com/stretchr/testify/assert"
)

func setupTestDB_PendingRequest(t *testing.T) {
	db.InitDB()
}

func TestPendingRequest(t *testing.T) {
	setupTestDB_PendingRequest(t)

	t.Run("正常系: id2 への保留中リクエスト", func(t *testing.T) {
		requests, err := PendingRequest("id2")
		assert.NoError(t, err)

		var gotIDs []string
		for _, r := range requests {
			gotIDs = append(gotIDs, r.User1ID)
		}

		assert.Contains(t, gotIDs, "id1")
		assert.NotContains(t, gotIDs, "id3")
		assert.NotContains(t, gotIDs, "id4")
		assert.NotContains(t, gotIDs, "id2")
	})
}
