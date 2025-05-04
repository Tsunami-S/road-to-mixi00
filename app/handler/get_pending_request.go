package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/db"
	"minimal_sns_app/model"
	"net/http"
)

func GetPendingRequests(c echo.Context) error {
	userID := c.QueryParam("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id is required"})
	}

	var requests []model.FriendRequest
	if err := db.DB.Where("user2_id = ? AND status = 'pending'", userID).Find(&requests).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch requests"})
	}

	return c.JSON(http.StatusOK, requests)
}
