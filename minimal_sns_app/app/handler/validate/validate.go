package validate

import (
	"github.com/labstack/echo/v4"
	"strconv"
	"net/http"
)

func ParseAndValidatePagination(c echo.Context) (limit int, page int, err error) {
	limitStr := c.QueryParam("limit")
	pageStr := c.QueryParam("page")

	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return 0, 0, echo.NewHTTPError(http.StatusBadRequest, "error: invalid limit")
	}

	page, err = strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		return 0, 0, echo.NewHTTPError(http.StatusBadRequest, "error: invalid page")
	}

	return limit, page, nil
}
