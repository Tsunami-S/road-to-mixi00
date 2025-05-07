package create

import (
	"github.com/labstack/echo/v4"
	repo_create "minimal_sns_app/repository/create"
	"net/http"
)

func AddNewUser(c echo.Context) error {
	id := c.QueryParam("id")
	name := c.QueryParam("name")

	// validation
	if id == "" || len(id) > 20 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id must have 1 ~ 20 characters"})
	}
	if name == "" || len(name) > 64 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name must have 1 ~ 64 characters"})
	}

	// check uniqueness
	exists, err := repo_create.IsUserIDExists(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to check user ID uniqueness"})
	}
	if exists {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user ID already exists"})
	}

	// create user
	if err := repo_create.User(id, name); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user added"})
}
