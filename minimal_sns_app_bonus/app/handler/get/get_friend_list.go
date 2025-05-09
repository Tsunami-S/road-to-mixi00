package get

import (
	"net/http"

	"minimal_sns_app/interfaces"

	"github.com/labstack/echo/v4"
)

type FriendHandler struct {
	Validator interfaces.Validator
	Repo      interfaces.FriendRepository
}

func NewFriendHandler(v interfaces.Validator, r interfaces.FriendRepository) *FriendHandler {
	return &FriendHandler{
		Validator: v,
		Repo:      r,
	}
}

func (h *FriendHandler) Friend(c echo.Context) error {
	userID := c.QueryParam("id")

	if userID == "" || len(userID) > 20 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "id must be a non-empty string up to 20 characters",
		})
	}

	exists, err := h.Validator.UserExists(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if !exists {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user not found"})
	}

	friends, err := h.Repo.GetFriends(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(friends) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no friends found"})
	}

	return c.JSON(http.StatusOK, friends)
}
