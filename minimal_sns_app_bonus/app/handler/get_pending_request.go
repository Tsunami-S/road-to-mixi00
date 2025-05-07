package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/repository"
	"net/http"
)

func GetPendingRequests(c echo.Context) error {
	userID := c.QueryParam("user_id")

	if valid, err := IsValidUserId(userID); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id: " + err.Error()})
	}

	requests, err := repository.GetPendingRequestsForUser(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch requests"})
	}

	if len(requests) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no pending requests found"})
	}

	return c.JSON(http.StatusOK, requests)
}
