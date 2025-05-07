package create

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/handler/validate"
	repo_create "minimal_sns_app/repository/create"
	"net/http"
)

func RespondRequest(c echo.Context) error {
	user1ID := c.QueryParam("user1_id")
	user2ID := c.QueryParam("user2_id")
	action := c.QueryParam("action")

	// validation
	if valid, err := validate.IsValidUserId(user1ID); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user1_id: " + err.Error()})
	}
	if valid, err := validate.IsValidUserId(user2ID); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user2_id: " + err.Error()})
	}
	if user1ID == user2ID {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user IDs"})
	}
	if action != "accepted" && action != "rejected" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid action"})
	}

	// check & update friend request
	req, err := repo_create.FindRequest(user1ID, user2ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "request not found or already handled"})
	}
	if err := repo_create.UpdateRequest(req, action); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update request"})
	}

	// create friend link if accepted
	if action == "accepted" {
		if err := repo_create.FriendLink(user1ID, user2ID); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create friendship"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "request " + action + "ed"})
}
