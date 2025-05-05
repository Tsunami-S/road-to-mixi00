package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/db"
	"minimal_sns_app/model"
	"net/http"
)

func AddBlockList(c echo.Context) error {
	user1ID := c.QueryParam("user1_id")
	user2ID := c.QueryParam("user2_id")

	// validation
	if len(user1ID) == 0 || len(user2ID) == 0 || len(user1ID) > 20 || len(user2ID) > 20 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user IDs"})
	}
	if user1ID == user2ID {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot block yourself"})
	}

	// check BLOCK status
	var existing model.BlockList
	if err := db.DB.Where("user1_id = ? AND user2_id = ?", user1ID, user2ID).First(&existing).Error; err == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "already blocked"})
	}

	// check FRIEND status
	if err := db.DB.Where(
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		user1ID, user2ID, user2ID, user1ID,
	).Delete(&model.FriendLink{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete friendship"})
	}

	// check REQUEST status
	if err := db.DB.Model(&model.FriendRequest{}).
		Where(
			"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
			user1ID, user2ID, user2ID, user1ID,
		).
		Where("status = ?", "pending").
		Update("status", "rejected").Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to reject friend request"})
	}

	// block user
	block := model.BlockList{User1ID: user1ID, User2ID: user2ID}
	if err := db.DB.Create(&block).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to block user"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user blocked and friendship removed if existed"})
}
