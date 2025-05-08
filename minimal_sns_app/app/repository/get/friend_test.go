package get

import (
	"testing"

	"minimal_sns_app/db"

	"github.com/stretchr/testify/assert"
)

func setupTestDB_FriendRepo(t *testing.T) {
	db.InitDB()
}

func TestFriendRepository_Friend(t *testing.T) {
	setupTestDB_FriendRepo(t)

	tests := []struct {
		name     string
		userID   int
		expected []string
	}{
		{
			name:     "フレンドあり（ブロックなし）",
			userID:   1,
			expected: []string{"user02", "user03", "user04", "user05", "user06", "user07", "user08", "user09", "user10"},
		},
		{
			name:     "ブロックされているユーザーは除外される",
			userID:   43,
			expected: []string{},
		},
		{
			name:     "存在しないユーザー",
			userID:   9999,
			expected: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Friend(tc.userID)
			assert.NoError(t, err)

			var actualNames []string
			for _, f := range result {
				actualNames = append(actualNames, f.Name)
			}

			assert.ElementsMatch(t, tc.expected, actualNames)
		})
	}
}
