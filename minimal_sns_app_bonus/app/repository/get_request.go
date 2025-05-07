package repository

import (
	"minimal_sns_app/db"
	"minimal_sns_app/model"
)

func GetPendingRequestsForUser(userID string) ([]model.FriendRequest, error) {
	var requests []model.FriendRequest

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

	err := db.DB.Raw(query, userID).Scan(&requests).Error
	return requests, err
}
