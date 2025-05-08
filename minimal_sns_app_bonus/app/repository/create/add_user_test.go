package create

import (
	"testing"

	"minimal_sns_app/db"
	"minimal_sns_app/model"

	"github.com/stretchr/testify/assert"
)

func setupTestDB_User(t *testing.T) {
	db.InitDB()
	db.DB.Exec("DELETE FROM users")
}

func TestIsUserIDExists(t *testing.T) {
	setupTestDB_User(t)

	db.DB.Create(&model.User{UserID: "user123", Name: "Test User"})

	t.Run("ユーザーが存在する", func(t *testing.T) {
		exists, err := IsUserIDExists("user123")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("ユーザーが存在しない", func(t *testing.T) {
		exists, err := IsUserIDExists("nonexistent")
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestCreateUser(t *testing.T) {
	setupTestDB_User(t)

	err := CreateUser("new_user", "New User")
	assert.NoError(t, err)

	var user model.User
	err = db.DB.Where("user_id = ?", "new_user").First(&user).Error
	assert.NoError(t, err)
	assert.Equal(t, "New User", user.Name)
}
