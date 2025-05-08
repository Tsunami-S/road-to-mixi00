package get

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"minimal_sns_app/interfaces"
)

type FriendOfFriendHandler struct {
	Validator interfaces.Validator
	Repo      interfaces.FriendOfFriendRepository
}

func NewFriendOfFriendHandler(v interfaces.Validator, r interfaces.FriendOfFriendRepository) *FriendOfFriendHandler {
	return &FriendOfFriendHandler{
		Validator: v,
		Repo:      r,
	}
}

func (h *FriendOfFriendHandler) FriendOfFriend(c echo.Context) error {
	userID := c.QueryParam("id")

	if userID == "" || len(userID) > 20 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id must be a non-empty string up to 20 characters"})
	}

	exists, err := h.Validator.UserExists(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if !exists {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user not found"})
	}

	result, err := h.Repo.GetFriendOfFriend(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(result) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no friends of friends found"})
	}

	return c.JSON(http.StatusOK, result)
}
