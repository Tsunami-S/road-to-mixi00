package get

import (
	"testing"

	"minimal_sns_app/db"

	"github.com/stretchr/testify/assert"
)

func setupTestDB_Friend(t *testing.T) {
	db.InitDB()
}

func TestFriend(t *testing.T) {
	setupTestDB_Friend(t)
	repo := &RealFriendRepository{}

	t.Run("正常系: id1 の友達取得", func(t *testing.T) {
		friends, err := repo.GetFriends("id1")
		assert.NoError(t, err)

		gotIDs := map[string]bool{}
		for _, f := range friends {
			gotIDs[f.ID] = true
		}

		assert.True(t, gotIDs["id2"], "id2 は含まれるべき")
		assert.True(t, gotIDs["id3"], "id3 は含まれるべき")
		assert.False(t, gotIDs["id4"], "id4 はブロックしているため除外されるべき")
		assert.False(t, gotIDs["id5"], "id5 からブロックされているため除外されるべき")
	})
}
