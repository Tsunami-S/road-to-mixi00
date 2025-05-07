package get_all

import (
	repo_all "minimal_sns_app/repository/get_all"
	"net/http"

	"github.com/labstack/echo/v4"
)

func FriendLinks(c echo.Context) error {
	links, err := repo_all.FriendLinks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch friend links"})
	}
	return c.JSON(http.StatusOK, links)
}
