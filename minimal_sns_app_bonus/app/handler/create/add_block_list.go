package create

import (
	"net/http"

	"minimal_sns_app/handler/validate"
	"minimal_sns_app/model"
	repo_create "minimal_sns_app/repository/create"

	"github.com/labstack/echo/v4"
)

func AddBlockList(c echo.Context) error {
	var req model.BlockRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	// validation
	if valid, err := validate.IsValidUserId(req.User1ID); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user1_id: " + err.Error()})
	}
	if valid, err := validate.IsValidUserId(req.User2ID); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user2_id: " + err.Error()})
	}
	if req.User1ID == req.User2ID {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot block yourself"})
	}

	// check already blocked
	blocked, err := repo_create.IsBlocked(req.User1ID, req.User2ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if blocked {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "already blocked"})
	}

	// remove friendship
	if err := repo_create.DeleteFriendLink(req.User1ID, req.User2ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete friendship"})
	}

	// reject pending requests
	if err := repo_create.RejectRequests(req.User1ID, req.User2ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to reject friend request"})
	}

	// add to block list
	if err := repo_create.Block(req.User1ID, req.User2ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to block user"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user blocked and friendship removed if existed"})
}
