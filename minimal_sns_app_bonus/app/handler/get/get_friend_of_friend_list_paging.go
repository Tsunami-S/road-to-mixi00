package get

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/handler/validation"
	repo_get "minimal_sns_app/repository/get"
	"net/http"
)

func FriendOfFriendPaging(c echo.Context) error {
	userID := c.QueryParam("id")

	// validation
	if valid, err := validation.IsValidUserId(userID); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id: " + err.Error()})
	}
	limit, page, err := validation.ParseAndValidatePagination(c)
	if err != nil {
		return err
	}

	// get friend list with paging
	offset := (page - 1) * limit
	result, err := repo_get.FriendOfFriendPaging(userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(result) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no friends of friends found"})
	}

	return c.JSON(http.StatusOK, result)
}
