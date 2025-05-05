package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/db"
	"minimal_sns_app/model"
	"net/http"
)

func RequestFriend(c echo.Context) error {
	user1ID := c.QueryParam("user1_id")
	user2ID := c.QueryParam("user2_id")

	// validation
	if len(user1ID) == 0 || len(user2ID) == 0 || user1ID == user2ID {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user IDs"})
	}

	// check FRIEND status
	var friend model.FriendLink
	err_friend := db.DB.Where(
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		user1ID, user2ID, user2ID, user1ID,
	).First(&friend).Error
	if err_friend == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "you are already friends"})
	}

	// check BLOCK status
	var block model.BlockList
	err_block := db.DB.Where(
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		user1ID, user2ID, user2ID, user1ID,
	).First(&block).Error
	if err_block == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot send friend request due to block"})
	}

	// check APPLICATION status
	var reverse model.FriendRequest
	if err := db.DB.Where("user1_id = ? AND user2_id = ? AND status = 'pending'", user2ID, user1ID).First(&reverse).Error; err == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "you already have a pending request from this user"})
	}
	var existing model.FriendRequest
	if err := db.DB.Where("user1_id = ? AND user2_id = ?", user1ID, user2ID).First(&existing).Error; err == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "friend request already sent"})
	}

	// send request
	req := model.FriendRequest{User1ID: user1ID, User2ID: user2ID, Status: "pending"}
	if err := db.DB.Create(&req).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to send friend request"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "friend request sent"})
}
