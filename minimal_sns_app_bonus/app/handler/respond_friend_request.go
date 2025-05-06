package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/db"
	"minimal_sns_app/model"
	"net/http"
)

func RespondFriendRequest(c echo.Context) error {
	user1ID := c.QueryParam("user1_id")
	user2ID := c.QueryParam("user2_id")
	action := c.QueryParam("action")

	// validation
	if valid, err := isValidUserId(user1ID); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user1_id: " + err.Error()})
	}
	if valid, err := isValidUserId(user2ID); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user2_id: " + err.Error()})
	}
	if user1ID == user2ID {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user IDs"})
	}
	if action != "accepted" && action != "rejected" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid action"})
	}

	// check REQUEST status
	var req model.FriendRequest
	err := db.DB.Where("user1_id = ? AND user2_id = ? AND status = 'pending'", user1ID, user2ID).First(&req).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "request not found or already handled"})
	}

	// update status
	if err := db.DB.Model(&req).Update("status", action).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update request"})
	}

	// add friend
	if action == "accepted" {
		link := model.FriendLink{User1ID: user1ID, User2ID: user2ID}
		if err := db.DB.Create(&link).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create friendship"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "request " + action + "ed"})
}
