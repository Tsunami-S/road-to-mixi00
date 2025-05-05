package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/db"
	"minimal_sns_app/model"
	"net/http"
)

func AddNewUser(c echo.Context) error {
	id := c.QueryParam("id")
	name := c.QueryParam("name")

	// validation
	if id == "" || len(id) > 20 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id must be 1 ~ 20 characters"})
	}
	if name == "" || len(name) > 64 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name must be 1 ~ 64 characters"})
	}

	// add new user
	user := model.User{UserID: id, Name: name}
	if err := db.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user added"})
}
