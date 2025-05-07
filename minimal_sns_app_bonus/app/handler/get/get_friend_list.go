package get

import (
	"minimal_sns_app/handler/validate"
	repo_get "minimal_sns_app/repository/get"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Friend(c echo.Context) error {
	userID := c.QueryParam("id")

	// validation
	if valid, err := validate.IsValidUserId(userID); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id: " + err.Error()})
	}

	// get friend list
	friends, err := repo_get.Friend(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(friends) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no friends found"})
	}

	return c.JSON(http.StatusOK, friends)
}
