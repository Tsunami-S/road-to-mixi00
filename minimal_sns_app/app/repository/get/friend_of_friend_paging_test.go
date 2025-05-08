package get

import (
	"testing"

	"minimal_sns_app/db"

	"github.com/stretchr/testify/assert"
)

func setupTestDB_FOFPaging(t *testing.T) {
	db.InitDB()
}

func TestFriendOfFriendPaging(t *testing.T) {
	setupTestDB_FOFPaging(t)

	tests := []struct {
		name     string
		userID   int
		limit    int
		offset   int
		expected []string
	}{
		{
			name:     "user1の友達の友達（ページ1）",
			userID:   1,
			limit:    2,
			offset:   0,
			expected: []string{"user11", "user12"},
		},
		{
			name:     "user1の友達の友達（ページ2）",
			userID:   1,
			limit:    2,
			offset:   2,
			expected: []string{"user13", "user31"},
		},
		{
			name:     "user33の友達の友達（多数blockで空）",
			userID:   33,
			limit:    2,
			offset:   0,
			expected: []string{},
		},
		{
			name:     "存在しないユーザー（ページ指定あり）",
			userID:   9999,
			limit:    2,
			offset:   0,
			expected: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := FriendOfFriendPaging(tc.userID, tc.limit, tc.offset)
			assert.NoError(t, err)

			var names []string
			for _, f := range result {
				names = append(names, f.Name)
			}

			assert.ElementsMatch(t, tc.expected, names)
		})
	}
}
