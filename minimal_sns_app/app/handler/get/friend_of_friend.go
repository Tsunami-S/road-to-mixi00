package get

import (
	"net/http"
	"strconv"

	"minimal_sns_app/interfaces"

	"github.com/labstack/echo/v4"
)

type FriendOfFriendHandler struct {
	Validator interfaces.UserValidator
	Repo      interfaces.FriendOfFriendRepository
}

func NewFriendOfFriendHandler(v interfaces.UserValidator, r interfaces.FriendOfFriendRepository) *FriendOfFriendHandler {
	return &FriendOfFriendHandler{Validator: v, Repo: r}
}

func (h *FriendOfFriendHandler) FriendOfFriend(c echo.Context) error {
	idStr := c.QueryParam("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id must be a positive integer"})
	}

	exist, err := h.Validator.UserExists(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if !exist {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user not found"})
	}

	result, err := h.Repo.GetFriendOfFriend(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(result) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no friends of friends found"})
	}

	return c.JSON(http.StatusOK, result)
}
