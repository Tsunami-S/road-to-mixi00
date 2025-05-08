package get

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"minimal_sns_app/interfaces"
)

type PendingRequestHandler struct {
	Validator interfaces.Validator
	Repo      interfaces.FriendRequestRepository
}

func NewPendingRequestHandler(v interfaces.Validator, r interfaces.FriendRequestRepository) *PendingRequestHandler {
	return &PendingRequestHandler{
		Validator: v,
		Repo:      r,
	}
}

func (h *PendingRequestHandler) PendingRequests(c echo.Context) error {
	userID := c.QueryParam("user_id")
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

	requests, err := h.Repo.GetPendingRequests(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(requests) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no pending requests found"})
	}

	return c.JSON(http.StatusOK, requests)
}
