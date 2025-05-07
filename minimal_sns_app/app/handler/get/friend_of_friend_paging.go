package get

import (
	"net/http"
	"strconv"

	handle_valid "minimal_sns_app/handler/validate"
	repo_get "minimal_sns_app/repository/get"
	rep_valid "minimal_sns_app/repository/validate"

	"github.com/labstack/echo/v4"
)

func FriendOfFriendPaging(c echo.Context) error {
	idStr := c.QueryParam("id")

	// validation
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id must be a positive integer"})
	}
	limit, page, err := handle_valid.ParseAndValidatePagination(c)
	if err != nil {
		return err
	}
	exist, err := rep_valid.UserExists(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if !exist {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user not found"})
	}

	// get friend list with paging
	offset := (page - 1) * limit
	result, err := repo_get.FriendOfFriendPaging(id, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(result) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no friends of friends found"})
	}

	return c.JSON(http.StatusOK, result)
}
