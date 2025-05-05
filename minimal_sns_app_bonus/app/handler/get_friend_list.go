package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/db"
	"minimal_sns_app/model"
	"net/http"
)

func GetFriendList(c echo.Context) error {
	id := c.QueryParam("id")

	// validation
	if id == "" || len(id) > 20 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id must be 1 ~ 20 characters"})
	}
	exist, err := userExists(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "DB error"})
	}
	if !exist {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user not found"})
	}

	// get friend list
	friends, err := getFriendsByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(friends) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no friends found"})
	}

	return c.JSON(http.StatusOK, friends)
}

func getFriendsByID(id string) ([]model.Friend, error) {
	var friends []model.Friend
	query := `
        SELECT u.user_id AS id, u.name
        FROM friend_link f
        JOIN users u ON u.user_id = f.user2_id
        WHERE f.user1_id = ?
        AND NOT EXISTS (
            SELECT 1 FROM block_list
            WHERE (user1_id = f.user1_id AND user2_id = f.user2_id)
               OR (user1_id = f.user2_id AND user2_id = f.user1_id)
        )
        UNION
        SELECT u.user_id AS id, u.name
        FROM friend_link f
        JOIN users u ON u.user_id = f.user1_id
        WHERE f.user2_id = ?
        AND NOT EXISTS (
            SELECT 1 FROM block_list
            WHERE (user1_id = f.user1_id AND user2_id = f.user2_id)
               OR (user1_id = f.user2_id AND user2_id = f.user1_id)
        );
    `

	err := db.DB.Raw(query, id, id).Scan(&friends).Error

	return friends, err
}
