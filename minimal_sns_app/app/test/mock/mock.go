package mock

import (
	"minimal_sns_app/model"

	"github.com/labstack/echo/v4"
)

// UserValidator mock
type UserValidatorMock struct {
	Exist bool
	Err   error
}

func (m *UserValidatorMock) UserExists(id int) (bool, error) {
	return m.Exist, m.Err
}

// FriendRepository mock
type FriendRepoMock struct {
	Result []model.Friend
	Err    error
}

func (m *FriendRepoMock) GetFriends(id int) ([]model.Friend, error) {
	return m.Result, m.Err
}

type FriendOfFriendRepoMock struct {
	Result []model.Friend
	Err    error
}

func (m *FriendOfFriendRepoMock) GetFriendOfFriend(id int) ([]model.Friend, error) {
	return m.Result, m.Err
}

type PaginationValidatorMock struct {
	Limit int
	Page  int
	Err   error
}

func (m *PaginationValidatorMock) ParseAndValidatePagination(c echo.Context) (int, int, error) {
	return m.Limit, m.Page, m.Err
}

type FriendOfFriendPagingRepoMock struct {
	Result []model.Friend
	Err    error
}

func (m *FriendOfFriendPagingRepoMock) GetFriendOfFriendByIDWithPaging(id, limit, offset int) ([]model.Friend, error) {
	return m.Result, m.Err
}
