package get_all

import (
	"net/http"

	repo_all "minimal_sns_app/repository/get_all"

	"github.com/labstack/echo/v4"
)

func RequestList(c echo.Context) error {
	requests, err := repo_all.FriendRequests()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, requests)
}
