package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/repository"
	"net/http"
)

func RequestFriend(c echo.Context) error {
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot request yourself"})
	}

	// check block
	blocked, err := repository.IsBlockedEachOther(user1ID, user2ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "DB error"})
	}
	if blocked {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot send friend request due to block"})
	}

	// check already friends
	alreadyFriends, err := repository.IsAlreadyFriends(user1ID, user2ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "DB error"})
	}
	if alreadyFriends {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "you are already friends"})
	}

	// check reverse request
	hasReverse, err := repository.HasPendingRequestFrom(user2ID, user1ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "DB error"})
	}
	if hasReverse {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "you already have a pending request from this user"})
	}

	// check same direction
	hasSent, err := repository.HasAlreadyRequested(user1ID, user2ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "DB error"})
	}
	if hasSent {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "friend request already sent"})
	}

	// create
	if err := repository.CreateFriendRequest(user1ID, user2ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to send friend request"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "friend request sent"})
}
