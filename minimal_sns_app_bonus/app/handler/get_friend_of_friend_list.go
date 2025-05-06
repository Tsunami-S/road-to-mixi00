package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/db"
	"minimal_sns_app/model"
	"net/http"
)

func GetFriendOfFriendList(c echo.Context) error {
	userID := c.QueryParam("id")

	// validation
	if valid, err := IsValidUserId(userID); !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id: " + err.Error()})
	}

	// get friend of friend list
	result, err := getFriendOfFriendByIDWithFilter(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(result) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no friends of friends found"})
	}

	return c.JSON(http.StatusOK, result)
}

func getFriendOfFriendByIDWithFilter(id string) ([]model.Friend, error) {
	var result []model.Friend

	query := `
	SELECT DISTINCT u.user_id AS id, u.name

	-- get DIRECT FRIENDS of the given user from friend_link
	FROM (
		SELECT CASE 
				 WHEN user1_id = ? THEN user2_id
				 WHEN user2_id = ? THEN user1_id
			   END AS friend_id
		FROM friend_link
		WHERE user1_id = ? OR user2_id = ?
	) AS direct

	-- get friends of those direct friends
	JOIN friend_link AS second
	  ON second.user1_id = direct.friend_id OR second.user2_id = direct.friend_id

	-- get user info of the friends_of_friends
	JOIN users u
	  ON u.user_id = IF(second.user1_id = direct.friend_id, second.user2_id, second.user1_id)

	-- exclude my id
	WHERE u.user_id != ? 

	-- exclude users who are already direct friends
	  AND u.user_id NOT IN (
		  SELECT CASE 
				   WHEN user1_id = ? THEN user2_id
				   ELSE user1_id
				 END
		  FROM friend_link
		  WHERE user1_id = ? OR user2_id = ?
	  )

	-- exclude users who are in block list with the original user
	  AND u.user_id NOT IN (
		  SELECT user2_id FROM block_list WHERE user1_id = ?
		  UNION
		  SELECT user1_id FROM block_list WHERE user2_id = ?
	  )

	-- exclude friends_of_friends if their connection is a block relationship
	  AND (second.user1_id, second.user2_id) NOT IN (
		  SELECT user1_id, user2_id FROM block_list
		  UNION
		  SELECT user2_id, user1_id FROM block_list
	  )
	`

	err := db.DB.Raw(query,
		id, id, // CASE: user1_id = ? || user2_id = ?
		id, id, // WHERE: user1_id = ? OR user2_id = ?
		id,         // WHERE u.user_id != ?
		id, id, id, // NOT IN: friend_link
		id, id, // NOT IN: block_list
	).Scan(&result).Error
	return result, err
}
