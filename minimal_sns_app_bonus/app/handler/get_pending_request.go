package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/db"
	"minimal_sns_app/model"
	"net/http"
)

func GetPendingRequests(c echo.Context) error {
	userID := c.QueryParam("user_id")

	if valid, err := IsValidUserId(userID); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id: " + err.Error()})
	}

	query := `
	SELECT fr.*
	FROM friend_requests fr
	WHERE fr.user2_id = ?
	  AND fr.status = 'pending'
	  AND fr.user1_id != fr.user2_id
	  AND NOT EXISTS (
	    SELECT 1 FROM block_list b
	    WHERE 
	      (b.user1_id = fr.user1_id AND b.user2_id = fr.user2_id)
	      OR (b.user1_id = fr.user2_id AND b.user2_id = fr.user1_id)
	  );
	`

	var requests []model.FriendRequest
	if err := db.DB.Raw(query, userID).Scan(&requests).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch requests"})
	}

	if len(requests) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no pending requests found"})
	}

	return c.JSON(http.StatusOK, requests)
}
