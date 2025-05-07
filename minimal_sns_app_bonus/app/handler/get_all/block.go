package get_all

import (
	"net/http"

	repo_all "minimal_sns_app/repository/get_all"

	"github.com/labstack/echo/v4"
)

func BlockList(c echo.Context) error {
	blocks, err := repo_all.BlockList()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch block list"})
	}
	return c.JSON(http.StatusOK, blocks)
}
