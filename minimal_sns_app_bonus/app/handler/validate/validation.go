package validate

import (
	"errors"
	"net/http"
	"strconv"

	repo_validate "minimal_sns_app/repository/validate"

	"github.com/labstack/echo/v4"
)

type RealValidator struct {
	UserExistsFunc func(id string) (bool, error)
}

func (v *RealValidator) UserExists(id string) (bool, error) {
	if id == "" || len(id) > 20 {
		return false, errors.New("invalid user ID format")
	}

	checkFunc := v.UserExistsFunc
	if checkFunc == nil {
		checkFunc = repo_validate.UserExists
	}

	exists, err := checkFunc(id)
	if err != nil {
		return false, errors.New("DB error while checking user ID")
	}
	if !exists {
		return false, errors.New("user not found")
	}
	return true, nil
}

type RealPaginationValidator struct{}

func (r *RealPaginationValidator) ParseAndValidatePagination(c echo.Context) (int, int, error) {
	limitStr := c.QueryParam("limit")
	pageStr := c.QueryParam("page")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return 0, 0, echo.NewHTTPError(http.StatusBadRequest, "error: invalid limit")
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		return 0, 0, echo.NewHTTPError(http.StatusBadRequest, "error: invalid page")
	}

	return limit, page, nil
}
