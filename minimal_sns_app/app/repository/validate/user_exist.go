package validate

import (
	"gorm.io/gorm"
	"net/http"
	"strconv"

	"minimal_sns_app/db"
	"minimal_sns_app/model"

	"github.com/labstack/echo/v4"
)

type RealValidator struct{}

func (r *RealValidator) UserExists(id int) (bool, error) {
	return UserExists(id)
}

func UserExists(id int) (bool, error) {
	var user model.User
	err := db.DB.Where("user_id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

type RealPaginationValidator struct{}

func (r *RealPaginationValidator) ParseAndValidatePagination(c echo.Context) (limit int, page int, err error) {
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
