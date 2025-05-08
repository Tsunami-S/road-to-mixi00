package get

import (
	"testing"

	"minimal_sns_app/db"

	"github.com/stretchr/testify/assert"
)

func TestFriendOfFriend(t *testing.T) {
	db.InitDB()
	repo := &RealFriendOfFriendRepository{}

	tests := []struct {
		name      string
		userID    string
		wantNames []string
		notWant   []string
		expectErr bool
	}{
		{
			name:      "正常系: id49の友達の友達",
			userID:    "id49",
			wantNames: []string{"user51"},
			expectErr: false,
		},
		{
			name:      "異常系: 存在しないユーザー",
			userID:    "nonexistent",
			wantNames: []string{},
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			results, err := repo.GetFriendOfFriend(tc.userID)

			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			names := map[string]bool{}
			for _, r := range results {
				names[r.Name] = true
			}

			for _, want := range tc.wantNames {
				assert.True(t, names[want], "want name %s to be present", want)
			}
			for _, notWant := range tc.notWant {
				assert.False(t, names[notWant], "name %s should not be present", notWant)
			}
		})
	}
}
