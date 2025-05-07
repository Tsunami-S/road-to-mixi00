package create

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/handler/validate"
	repo_create "minimal_sns_app/repository/create"
	"net/http"
)

func AddBlockList(c echo.Context) error {
	user1ID := c.QueryParam("user1_id")
	user2ID := c.QueryParam("user2_id")

	// validation
	if valid, err := validate.IsValidUserId(user1ID); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user1_id: " + err.Error()})
	}
	if valid, err := validate.IsValidUserId(user2ID); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user2_id: " + err.Error()})
	}
	if user1ID == user2ID {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot block yourself"})
	}

	// check already blocked
	blocked, err := repo_create.IsBlocked(user1ID, user2ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if blocked {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "already blocked"})
	}

	// remove friendship
	if err := repo_create.DeleteFriendLink(user1ID, user2ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete friendship"})
	}

	// reject pending requests
	if err := repo_create.RejectRequests(user1ID, user2ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to reject friend request"})
	}

	// add to block list
	if err := repo_create.Block(user1ID, user2ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to block user"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user blocked and friendship removed if existed"})
}
