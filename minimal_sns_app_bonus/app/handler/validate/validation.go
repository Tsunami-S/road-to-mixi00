package validate

import (
	"errors"
	"net/http"
	"strconv"

	repo_valid "minimal_sns_app/repository/validate"

	"github.com/labstack/echo/v4"
)

func IsValidUserId(id string) (bool, error) {
	if id == "" || len(id) > 20 {
		return false, errors.New("invalid user ID format")
	}

	exists, err := repo_valid.UserExists(id)
	if err != nil {
		return false, errors.New("DB error while checking user ID")
	}
	if !exists {
		return false, errors.New("user ID not found")
	}

	return true, nil
}

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
