package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"minimal_sns_app/db"
	"minimal_sns_app/model"
	"strconv"
)

func isValidUserId(id string) (bool, error) {
	if id == "" || len(id) > 20 {
		return false, errors.New("invalid user ID format")
	}

	var count int64
	err := db.DB.Model(&model.User{}).Where("user_id = ?", id).Count(&count).Error
	if err != nil {
		return false, errors.New("DB error while checking user ID")
	}
	if count == 0 {
		return false, errors.New("user ID not found")
	}

	return true, nil
}

func userExists(id string) (bool, error) {
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

func parseAndValidatePagination(c echo.Context) (limit int, page int, err error) {
	limitStr := c.QueryParam("limit")
	pageStr := c.QueryParam("page")

	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return 0, 0, echo.NewHTTPError(400, "error: invalid limit")
	}

	page, err = strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		return 0, 0, echo.NewHTTPError(400, "error: invalid page")
	}

	return limit, page, nil
}
