package get

import (
	"testing"

	"minimal_sns_app/db"

	"github.com/stretchr/testify/assert"
)

func setupTestDB_FriendOfFriendPaging(t *testing.T) {
	db.InitDB()
}

func TestFriendOfFriendPaging(t *testing.T) {
	setupTestDB_FriendOfFriendPaging(t)

	repo := &RealFriendOfFriendPagingRepository{}

	t.Run("ページング付き友達の友達取得", func(t *testing.T) {
		result, err := repo.GetFriendOfFriendByIDWithPaging("id1", 2, 0)
		assert.NoError(t, err)
		assert.LessOrEqual(t, len(result), 2)

		result2, err := repo.GetFriendOfFriendByIDWithPaging("id1", 2, 2)
		assert.NoError(t, err)
		assert.LessOrEqual(t, len(result2), 2)

		for _, f := range append(result, result2...) {
			assert.NotEqual(t, "id5", f.ID)
		}
	})
}
