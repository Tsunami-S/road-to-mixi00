package get

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/repository/validate"
	repo_get "minimal_sns_app/repository/get"
	"net/http"
	"strconv"
)

func FriendOfFriend(c echo.Context) error {
	idStr := c.QueryParam("id")

	// validation
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id must be a positive integer"})
	}
	exist, err := validate.UserExists(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if !exist {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user not found"})
	}

	// get friend of friend list
	result, err := repo_get.FriendOfFriend(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(result) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no friends of friends found"})
	}

	return c.JSON(http.StatusOK, result)
}
