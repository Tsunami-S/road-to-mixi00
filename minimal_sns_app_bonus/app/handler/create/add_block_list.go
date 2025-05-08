package create

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"minimal_sns_app/interfaces"
	"minimal_sns_app/model"
)

type BlockHandler struct {
	Validator interfaces.Validator
	Repo      interfaces.BlockRepository
}

func NewBlockHandler(v interfaces.Validator, r interfaces.BlockRepository) *BlockHandler {
	return &BlockHandler{
		Validator: v,
		Repo:      r,
	}
}

func (h *BlockHandler) AddBlockList(c echo.Context) error {
	var req model.BlockRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	if req.User1ID == "" || req.User2ID == "" || len(req.User1ID) > 20 || len(req.User2ID) > 20 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user IDs"})
	}
	if req.User1ID == req.User2ID {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot block yourself"})
	}

	if exist, err := h.Validator.UserExists(req.User1ID); err != nil || !exist {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if exist, err := h.Validator.UserExists(req.User2ID); err != nil || !exist {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	blocked, err := h.Repo.IsBlocked(req.User1ID, req.User2ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if blocked {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "already blocked"})
	}

	if err := h.Repo.DeleteFriendLink(req.User1ID, req.User2ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete friendship"})
	}
	if err := h.Repo.RejectRequests(req.User1ID, req.User2ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to reject friend request"})
	}
	if err := h.Repo.Block(req.User1ID, req.User2ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to block user"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user blocked and friendship removed if existed"})
}
