package mock

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/model"
)

// Validator モック
type Validator struct {
	Exists bool
	Err    error
}

func (m *Validator) UserExists(id string) (bool, error) {
	return m.Exists, m.Err
}

// PaginationValidator モック
type PaginationValidator struct {
	Limit, Page int
	Err         error
}

func (m *PaginationValidator) ParseAndValidatePagination(c echo.Context) (int, int, error) {
	return m.Limit, m.Page, m.Err
}

// FriendRepo モック
type FriendRepo struct {
	Result []model.Friend
	Err    error
}

func (m *FriendRepo) GetFriends(id string) ([]model.Friend, error) {
	return m.Result, m.Err
}

// FriendOfFriendRepo モック
type FriendOfFriendRepo struct {
	Result []model.Friend
	Err    error
}

func (m *FriendOfFriendRepo) GetFriendOfFriend(id string) ([]model.Friend, error) {
	return m.Result, m.Err
}

// FriendOfFriendPagingRepo モック
type FriendOfFriendPagingRepo struct {
	Result []model.Friend
	Err    error
}

func (m *FriendOfFriendPagingRepo) GetFriendOfFriendByIDWithPaging(id string, limit, offset int) ([]model.Friend, error) {
	return m.Result, m.Err
}

// FriendRequestRepo モック
type FriendRequestRepo struct {
	Requests              []model.FriendRequest
	Err                   error
	HasRequestedResult    bool
	HasRequestedErr       error
	RequestResult         error
	IsBlockedResult       bool
	IsBlockedErr          error
	IsAlreadyFriendsValue bool
	IsAlreadyFriendsErr   error
	HasReverseRequest     bool
	HasReverseRequestErr  error
}

func (m *FriendRequestRepo) GetPendingRequests(userID string) ([]model.FriendRequest, error) {
	return m.Requests, m.Err
}

func (m *FriendRequestRepo) HasAlreadyRequested(user1, user2 string) (bool, error) {
	return m.HasRequestedResult, m.HasRequestedErr
}

func (m *FriendRequestRepo) Request(user1, user2 string) error {
	return m.RequestResult
}

func (m *FriendRequestRepo) IsBlockedEachOther(user1, user2 string) (bool, error) {
	return m.IsBlockedResult, m.IsBlockedErr
}

func (m *FriendRequestRepo) IsAlreadyFriends(user1, user2 string) (bool, error) {
	return m.IsAlreadyFriendsValue, m.IsAlreadyFriendsErr
}

func (m *FriendRequestRepo) HasPendingRequest(user1, user2 string) (bool, error) {
	return m.HasReverseRequest, m.HasReverseRequestErr
}

func (m *FriendRequestRepo) FindRequest(user1, user2 string) (*model.FriendRequest, error) {
	if len(m.Requests) == 0 {
		return nil, m.Err
	}
	return &m.Requests[0], m.Err
}

func (m *FriendRequestRepo) UpdateRequest(req *model.FriendRequest, status string) error {
	return m.Err
}

func (m *FriendRequestRepo) FriendLink(user1, user2 string) error {
	return m.Err
}

// UserRepository モック
type UserRepository struct {
	Exists bool
	Err    error
}

func (m *UserRepository) IsUserIDExists(id string) (bool, error) {
	return m.Exists, m.Err
}

func (m *UserRepository) CreateUser(id, name string) error {
	return m.Err
}

// BlockRepository モック
type BlockRepository struct {
	IsBlockedResult bool
	IsBlockedErr    error
	OpErr           error
}

func (m *BlockRepository) IsBlocked(user1, user2 string) (bool, error) {
	return m.IsBlockedResult, m.IsBlockedErr
}

func (m *BlockRepository) DeleteFriendLink(user1, user2 string) error {
	return m.OpErr
}

func (m *BlockRepository) RejectRequests(user1, user2 string) error {
	return m.OpErr
}

func (m *BlockRepository) Block(user1, user2 string) error {
	return m.OpErr
}

type MockRespondRepo struct{}

func (m *MockRespondRepo) FindRequest(user1, user2 string) (*model.FriendRequest, error) {
	return &model.FriendRequest{}, nil
}

func (m *MockRespondRepo) UpdateRequest(req *model.FriendRequest, action string) error {
	return nil
}

func (m *MockRespondRepo) CreateFriendLink(user1, user2 string) error {
	return nil
}

func (m *MockRespondRepo) RespondRequest(fromID, toID, action string) error {
	return nil
}
