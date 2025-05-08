package create

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"minimal_sns_app/interfaces"
	"minimal_sns_app/model"
)

type UserHandler struct {
	Repo interfaces.UserRepository
}

func NewUserHandler(r interfaces.UserRepository) *UserHandler {
	return &UserHandler{
		Repo: r,
	}
}

func (h *UserHandler) AddNewUser(c echo.Context) error {
	var req model.AddUserRequest

	// parse JSON body into struct
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	// validate ID and Name
	if req.ID == "" || len(req.ID) > 20 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id must have 1 ~ 20 characters"})
	}
	if req.Name == "" || len(req.Name) > 64 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name must have 1 ~ 64 characters"})
	}

	// check if ID already exists
	exists, err := h.Repo.IsUserIDExists(req.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if exists {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user ID already exists"})
	}

	// insert new user
	if err := h.Repo.CreateUser(req.ID, req.Name); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user added"})
}
