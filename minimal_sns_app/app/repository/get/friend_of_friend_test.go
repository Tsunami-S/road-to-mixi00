package get

import (
	"testing"

	"minimal_sns_app/db"

	"github.com/stretchr/testify/assert"
)

func setupTestDB_FOF(t *testing.T) {
	db.InitDB()
}

func TestFriendOfFriend(t *testing.T) {
	setupTestDB_FOF(t)

	tests := []struct {
		name     string
		userID   int
		expected []string
	}{
		{
			name:     "user1の友達の友達（直接の友達とブロック関係を除外）",
			userID:   1,
			expected: []string{"user11", "user12", "user13", "user31", "user32", "user36"},
		},
		{
			name:     "user33の友達の友達（block関係で結果が除外される）",
			userID:   33,
			expected: []string{},
		},
		{
			name:     "存在しないユーザーの結果は空",
			userID:   9999,
			expected: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := FriendOfFriend(tc.userID)
			assert.NoError(t, err)

			var names []string
			for _, f := range result {
				names = append(names, f.Name)
			}

			assert.ElementsMatch(t, tc.expected, names)
		})
	}
}
