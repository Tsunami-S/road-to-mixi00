package create

import (
	"net/http"

	"minimal_sns_app/interfaces"
	"minimal_sns_app/model"

	"github.com/labstack/echo/v4"
)

type FriendRespondHandler struct {
	Validator interfaces.Validator
	Repo      interfaces.FriendRespondRepository
}

func NewFriendRespondHandler(v interfaces.Validator, r interfaces.FriendRespondRepository) *FriendRespondHandler {
	return &FriendRespondHandler{
		Validator: v,
		Repo:      r,
	}
}

func (h *FriendRespondHandler) RespondRequest(c echo.Context) error {
	var req model.RespondRequestInput

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	if req.User1ID == "" || req.User2ID == "" || len(req.User1ID) > 20 || len(req.User2ID) > 20 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user IDs"})
	}
	if req.User1ID == req.User2ID {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user IDs"})
	}
	if req.Action != "accepted" && req.Action != "rejected" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid action"})
	}

	exists, err := h.Validator.UserExists(req.User1ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if !exists {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user ID not found"})
	}

	exists, err = h.Validator.UserExists(req.User2ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if !exists {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user ID not found"})
	}

	if req.User1ID == req.User2ID {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot request yourself"})
	}

	request, err := h.Repo.FindRequest(req.User1ID, req.User2ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "request not found or already handled"})
	}
	if err := h.Repo.UpdateRequest(request, req.Action); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update request"})
	}

	if req.Action == "accepted" {
		if err := h.Repo.CreateFriendLink(req.User1ID, req.User2ID); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create friendship"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "request " + req.Action + "ed"})
}
