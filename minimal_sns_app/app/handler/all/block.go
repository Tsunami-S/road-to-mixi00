package all

import (
	"github.com/labstack/echo/v4"
	repo_all "minimal_sns_app/repository/all"
	"net/http"
)

func BlockList(c echo.Context) error {
	blocks, err := repo_all.BlockList()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch block list"})
	}
	return c.JSON(http.StatusOK, blocks)
}

