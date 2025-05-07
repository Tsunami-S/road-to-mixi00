package create

import (
	"net/http"

	"minimal_sns_app/handler/validate"
	"minimal_sns_app/model"
	repo_create "minimal_sns_app/repository/create"

	"github.com/labstack/echo/v4"
)

func RequestFriend(c echo.Context) error {
	var req model.FriendRequestInput

	// bind json to struct
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot request yourself"})
	}

	// check block
	blocked, err := repo_create.IsBlockedEachOther(req.User1ID, req.User2ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if blocked {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot send friend request due to block"})
	}

	// check already friends
	alreadyFriends, err := repo_create.IsAlreadyFriends(req.User1ID, req.User2ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "DB error"})
	}
	if alreadyFriends {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "you are already friends"})
	}

	// check reverse request
	hasReverse, err := repo_create.HasPendingRequest(req.User2ID, req.User1ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "DB error"})
	}
	if hasReverse {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "you already have a pending request from this user"})
	}

	// check same direction
	hasSent, err := repo_create.HasAlreadyRequested(req.User1ID, req.User2ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "DB error"})
	}
	if hasSent {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "friend request already sent"})
	}

	// request
	if err := repo_create.Request(req.User1ID, req.User2ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to send friend request"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "friend request sent"})
}
