package create

import (
	"minimal_sns_app/model"
	repo_create "minimal_sns_app/repository/create"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AddNewUser(c echo.Context) error {
	var req model.AddUserRequest

	// bind json to struct
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	// validation
	if req.ID == "" || len(req.ID) > 20 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id must have 1 ~ 20 characters"})
	}
	if req.Name == "" || len(req.Name) > 64 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name must have 1 ~ 64 characters"})
	}
	exists, err := repo_create.IsUserIDExists(req.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to check user ID uniqueness"})
	}
	if exists {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user ID already exists"})
	}

	// add new user
	if err := repo_create.User(req.ID, req.Name); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user added"})
}
