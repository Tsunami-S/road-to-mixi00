package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/db"
	"minimal_sns_app/model"
	"net/http"
)

func GetFriendList(c echo.Context) error {
	id, err := parseAndValidateID(c)
	if err != nil {
		return err
	}

	exist, err := userExists(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "DB error"})
	}
	if !exist {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user not found"})
	}

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
        UNION
        SELECT u.user_id , u.name
        FROM friend_link f
        JOIN users u ON u.user_id = f.user1_id
        WHERE f.user2_id = ?;
    `

	err := db.DB.Raw(query, id, id).Scan(&friends).Error

	return friends, err
}
