package get

import (
	"minimal_sns_app/handler/validate"
	repo_get "minimal_sns_app/repository/get"
	"net/http"

	"github.com/labstack/echo/v4"
)

func FriendOfFriendPaging(c echo.Context) error {
	userID := c.QueryParam("id")

	// validation
	if valid, err := validate.IsValidUserId(userID); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id: " + err.Error()})
	}
	limit, page, err := validate.ParseAndValidatePagination(c)
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
