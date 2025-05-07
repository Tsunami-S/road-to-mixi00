package get

import (
	"net/http"

	"minimal_sns_app/handler/validate"
	repo_get "minimal_sns_app/repository/get"

	"github.com/labstack/echo/v4"
)

func FriendOfFriend(c echo.Context) error {
	userID := c.QueryParam("id")

	// validation
	if valid, err := validate.IsValidUserId(userID); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id: " + err.Error()})
	}

	// get friend of friend list
	result, err := repo_get.FriendOfFriend(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(result) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no friends of friends found"})
	}

	return c.JSON(http.StatusOK, result)
}
