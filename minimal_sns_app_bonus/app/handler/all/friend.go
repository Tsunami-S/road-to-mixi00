package all

import (
	"github.com/labstack/echo/v4"
	repo_all "minimal_sns_app/repository/all"
	"net/http"
)

func FriendLinks(c echo.Context) error {
	links, err := repo_all.FriendLinks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch friend links"})
	}
	return c.JSON(http.StatusOK, links)
}
