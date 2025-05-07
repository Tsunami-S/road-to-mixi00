package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/repository"
	"net/http"
)

func GetFriendOfFriendListPaging(c echo.Context) error {
	userID := c.QueryParam("id")

	// validation
	if valid, err := IsValidUserId(userID); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id: " + err.Error()})
	}
	limit, page, err := ParseAndValidatePagination(c)
	if err != nil {
		return err
	}

	// get friend list with paging
	offset := (page - 1) * limit
	result, err := repository.GetFriendOfFriendByIDWithPaging(userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(result) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no friends of friends found"})
	}

	return c.JSON(http.StatusOK, result)
}
