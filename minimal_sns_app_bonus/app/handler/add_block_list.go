package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/repository"
	"net/http"
)

func AddBlockList(c echo.Context) error {
	user1ID := c.QueryParam("user1_id")
	user2ID := c.QueryParam("user2_id")

	// validation
	if valid, err := IsValidUserId(user1ID); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user1_id: " + err.Error()})
	}
	if valid, err := IsValidUserId(user2ID); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user2_id: " + err.Error()})
	}
	if user1ID == user2ID {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot block yourself"})
	}

	// check already blocked
	blocked, err := repository.IsBlocked(user1ID, user2ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "DB error"})
	}
	if blocked {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "already blocked"})
	}

	// remove friendship
	if err := repository.DeleteFriendLink(user1ID, user2ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete friendship"})
	}

	// reject pending requests
	if err := repository.RejectPendingRequests(user1ID, user2ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to reject friend request"})
	}

	// add to block list
	if err := repository.CreateBlock(user1ID, user2ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to block user"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user blocked and friendship removed if existed"})
}
