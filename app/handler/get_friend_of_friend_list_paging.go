package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/db"
	"minimal_sns_app/model"
	"net/http"
)

func GetFriendOfFriendListPaging(c echo.Context) error {
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

	limit, page, err := parseAndValidatePagination(c)
	if err != nil {
		return err
	}

	offset := (page - 1) * limit
	result, err := getFriendOfFriendByIDWithPaging(id, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(result) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no friends of friends found"})
	}

	return c.JSON(http.StatusOK, result)
}

func getFriendOfFriendByIDWithPaging(id string, limit, offset int) ([]model.Friend, error) {
	var result []model.Friend

	query := `
	SELECT DISTINCT u.user_id AS id, u.name
	FROM (
		SELECT CASE 
				 WHEN user1_id = ? THEN user2_id
				 WHEN user2_id = ? THEN user1_id
			   END AS friend_id
		FROM friend_link
		WHERE user1_id = ? OR user2_id = ?
	) AS direct
	JOIN friend_link AS second
	  ON second.user1_id = direct.friend_id OR second.user2_id = direct.friend_id
	JOIN users u
	  ON u.user_id = IF(second.user1_id = direct.friend_id, second.user2_id, second.user1_id)
	WHERE u.user_id != ?
	  AND u.user_id NOT IN (
		  SELECT CASE 
				   WHEN user1_id = ? THEN user2_id
				   ELSE user1_id
				 END
		  FROM friend_link
		  WHERE user1_id = ? OR user2_id = ?
	  )
	  AND u.user_id NOT IN (
		  SELECT user2_id FROM block_list WHERE user1_id = ?
		  UNION
		  SELECT user1_id FROM block_list WHERE user2_id = ?
	  )
	LIMIT ? OFFSET ?
	`

	err := db.DB.Raw(query,
		id, id, id, id, // direct
		id,         // u.user_id != ?
		id, id, id, // not already friend
		id, id, // block list
		limit, offset,
	).Scan(&result).Error

	return result, err
}
