package create

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"minimal_sns_app/interfaces"
	"minimal_sns_app/model"
)

type RequestFriendHandler struct {
	Validator interfaces.Validator
	Repo      interfaces.FriendRequestRepository
}

func NewRequestFriendHandler(v interfaces.Validator, r interfaces.FriendRequestRepository) *RequestFriendHandler {
	return &RequestFriendHandler{
		Validator: v,
		Repo:      r,
	}
}

func (h *RequestFriendHandler) RequestFriend(c echo.Context) error {
	var req model.FriendRequestInput

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	exists, err := h.Validator.UserExists(req.User1ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if !exists {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user1_id: user ID not found"})
	}

	exists, err = h.Validator.UserExists(req.User2ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if !exists {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user2_id: user ID not found"})
	}

	if req.User1ID == req.User2ID {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot request yourself"})
	}

	blocked, err := h.Repo.IsBlockedEachOther(req.User1ID, req.User2ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if blocked {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot send friend request due to block"})
	}

	alreadyFriends, err := h.Repo.IsAlreadyFriends(req.User1ID, req.User2ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if alreadyFriends {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "you are already friends"})
	}

	hasReverse, err := h.Repo.HasPendingRequest(req.User2ID, req.User1ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "DB error"})
	}
	if hasReverse {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "you already have a pending request from this user"})
	}

	hasSent, err := h.Repo.HasPendingRequest(req.User1ID, req.User2ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if hasSent {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "friend request already sent"})
	}

	if err := h.Repo.Request(req.User1ID, req.User2ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to send friend request"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "friend request sent"})
}
