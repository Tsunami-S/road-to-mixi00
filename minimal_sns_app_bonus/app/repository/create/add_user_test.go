package create

import (
	"testing"

	"minimal_sns_app/db"
	"minimal_sns_app/model"

	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) {
	db.InitDB()
}

func TestIsUserIDExists(t *testing.T) {
	setupTestDB(t)
	repo := &RealUserRepository{}

	user := model.User{UserID: "existing_id", Name: "Test"}
	err := db.DB.Create(&user).Error
	assert.NoError(t, err)

	t.Run("存在するIDを確認", func(t *testing.T) {
		exists, err := repo.IsUserIDExists("existing_id")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("存在しないIDを確認", func(t *testing.T) {
		exists, err := repo.IsUserIDExists("nonexistent_id")
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestCreateUser(t *testing.T) {
	setupTestDB(t)
	repo := &RealUserRepository{}

	t.Run("正常にユーザーを追加できる", func(t *testing.T) {
		err := repo.CreateUser("new_user", "新しいユーザー")
		assert.NoError(t, err)

		var user model.User
		err = db.DB.Where("user_id = ?", "new_user").First(&user).Error
		assert.NoError(t, err)
		assert.Equal(t, "新しいユーザー", user.Name)
	})
}
