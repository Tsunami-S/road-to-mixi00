package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/repository"
	"net/http"
)

func GetFriendList(c echo.Context) error {
	userID := c.QueryParam("id")

	// validation
	if valid, err := IsValidUserId(userID); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id: " + err.Error()})
	}

	// get friend list
	friends, err := repository.GetFriendsByID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(friends) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no friends found"})
	}

	return c.JSON(http.StatusOK, friends)
}
