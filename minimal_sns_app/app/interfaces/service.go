package interfaces

import (
	"minimal_sns_app/model"

	"github.com/labstack/echo/v4"
)

type UserValidator interface {
	UserExists(id int) (bool, error)
}

type FriendRepository interface {
	GetFriends(id int) ([]model.Friend, error)
}

type FriendOfFriendRepository interface {
	GetFriendOfFriend(id int) ([]model.Friend, error)
}

type PaginationValidator interface {
	ParseAndValidatePagination(c echo.Context) (limit, page int, err error)
}

type FriendOfFriendPagingRepository interface {
	GetFriendOfFriendByIDWithPaging(id, limit, offset int) ([]model.Friend, error)
}
