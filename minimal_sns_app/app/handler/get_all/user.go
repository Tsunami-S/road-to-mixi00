package get_all

import (
	repo_all "minimal_sns_app/repository/get_all"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Users(c echo.Context) error {
	users, err := repo_all.Users()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch users"})
	}
	return c.JSON(http.StatusOK, users)
}
