package get

import (
	"strconv"
	"github.com/labstack/echo/v4"
	"minimal_sns_app/repository/validation"
	repo_get "minimal_sns_app/repository/get"
	"net/http"
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
	limit, page, err := validation.ParseAndValidatePagination(c)
	if err != nil {
		return err
	}
	exist, err := validation.UserExists(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "DB error"})
	}
	if !exist {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user not found"})
	}

	// get friend list with paging
	offset := (page - 1) * limit
	result, err := repo_get.GetFriendOfFriendByIDWithPaging(id, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(result) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no friends of friends found"})
	}

	return c.JSON(http.StatusOK, result)
}
