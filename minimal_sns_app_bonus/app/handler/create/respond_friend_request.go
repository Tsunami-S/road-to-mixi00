package create

import (
	"net/http"

	"minimal_sns_app/handler/validate"
	"minimal_sns_app/model"
	repo_create "minimal_sns_app/repository/create"

	"github.com/labstack/echo/v4"
)

func RespondRequest(c echo.Context) error {
	var req model.RespondRequestInput

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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user IDs"})
	}
	if req.Action != "accepted" && req.Action != "rejected" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid action"})
	}

	// check & update friend request
	request, err := repo_create.FindRequest(req.User1ID, req.User2ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "request not found or already handled"})
	}
	if err := repo_create.UpdateRequest(request, req.Action); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update request"})
	}

	// add link
	if req.Action == "accepted" {
		if err := repo_create.FriendLink(req.User1ID, req.User2ID); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create friendship"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "request " + req.Action + "ed"})
}
