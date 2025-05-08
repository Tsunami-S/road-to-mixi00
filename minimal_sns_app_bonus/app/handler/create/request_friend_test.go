package create

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"minimal_sns_app/interfaces"
	"minimal_sns_app/model"
)

type FriendRequestHandler struct {
	Validator interfaces.Validator
	Repo      interfaces.FriendRequestRepository
}

func NewFriendRequestHandler(v interfaces.Validator, r interfaces.FriendRequestRepository) *FriendRequestHandler {
	return &FriendRequestHandler{
		Validator: v,
		Repo:      r,
	}
}

func (h *FriendRequestHandler) RequestFriend(c echo.Context) error {
	var req model.FriendRequestInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	if req.User1ID == "" || req.User2ID == "" || len(req.User1ID) > 20 || len(req.User2ID) > 20 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user IDs"})
	}
	if req.User1ID == req.User2ID {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot request yourself"})
	}

	if blocked, err := h.Repo.IsBlockedEachOther(req.User1ID, req.User2ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	} else if blocked {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot send friend request due to block"})
	}

	if isFriend, err := h.Repo.IsAlreadyFriends(req.User1ID, req.User2ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	} else if isFriend {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "you are already friends"})
	}

	if hasReverse, err := h.Repo.HasPendingRequest(req.User2ID, req.User1ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "DB error"})
	} else if hasReverse {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "you already have a pending request from this user"})
	}

	if hasSent, err := h.Repo.HasAlreadyRequested(req.User1ID, req.User2ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	} else if hasSent {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "friend request already sent"})
	}

	if err := h.Repo.Request(req.User1ID, req.User2ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to send friend request"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "friend request sent"})
}
