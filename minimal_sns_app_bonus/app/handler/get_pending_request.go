package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/db"
	"minimal_sns_app/model"
	"net/http"
)

func GetPendingRequests(c echo.Context) error {
	userID := c.QueryParam("user_id")

	// validation
	if valid, err := isValidUserId(userID); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id: " + err.Error()})
	}

	// get request
	var requests []model.FriendRequest
	if err := db.DB.Where("user2_id = ? AND status = 'pending'", userID).Find(&requests).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch requests"})
	}

	return c.JSON(http.StatusOK, requests)
}
