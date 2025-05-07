package get

import (
	"github.com/labstack/echo/v4"
	repo_get "minimal_sns_app/repository/get"
	"minimal_sns_app/repository/validation"
	"net/http"
	"strconv"
)

func Friend(c echo.Context) error {
	idStr := c.QueryParam("id")

	// validation
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 || len(idStr) > 11 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id must be 0 ~ 99999999999"})
	}
	exist, err := validation.UserExists(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "DB error"})
	}
	if !exist {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user not found"})
	}

	// get friend list
	friends, err := repo_get.GetFriendsByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(friends) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no friends found"})
	}

	return c.JSON(http.StatusOK, friends)
}
