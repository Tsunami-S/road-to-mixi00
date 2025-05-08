package get

import (
	"minimal_sns_app/db"
	"minimal_sns_app/model"
)

type RealFriendOfFriendPagingRepository struct{}

func (r *RealFriendOfFriendPagingRepository) GetFriendOfFriendByIDWithPaging(id string, limit, offset int) ([]model.Friend, error) {
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
