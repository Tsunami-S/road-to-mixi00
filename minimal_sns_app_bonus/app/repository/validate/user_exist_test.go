package validate

import (
	"testing"

	"minimal_sns_app/db"
	"minimal_sns_app/model"

	"github.com/stretchr/testify/assert"
)

func setupTestDB_UserExists(t *testing.T) {
	db.InitDB()

	db.DB.Create(&model.User{UserID: "exist_user", Name: "テストユーザー"})
}

func TestUserExists(t *testing.T) {
	setupTestDB_UserExists(t)

	tests := []struct {
		name      string
		inputID   string
		wantExist bool
		wantErr   bool
	}{
		{
			name:      "存在するユーザーID",
			inputID:   "exist_user",
			wantExist: true,
			wantErr:   false,
		},
		{
			name:      "存在しないユーザーID",
			inputID:   "no_user",
			wantExist: false,
			wantErr:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			exist, err := UserExists(tc.inputID)
			assert.Equal(t, tc.wantExist, exist)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
