package get

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"minimal_sns_app/interfaces"
)

type FriendHandler struct {
	Validator interfaces.UserValidator
	Repo      interfaces.FriendRepository
}

func NewFriendHandler(v interfaces.UserValidator, r interfaces.FriendRepository) *FriendHandler {
	return &FriendHandler{Validator: v, Repo: r}
}

func (h *FriendHandler) Friend(c echo.Context) error {
	idStr := c.QueryParam("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 || len(idStr) > 11 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id must be 0 ~ 99999999999"})
	}
	exist, err := h.Validator.UserExists(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if !exist {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user not found"})
	}

	friends, err := h.Repo.GetFriends(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(friends) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no friends found"})
	}

	return c.JSON(http.StatusOK, friends)
}
