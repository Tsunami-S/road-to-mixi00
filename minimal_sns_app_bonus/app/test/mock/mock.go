package mock

import (
	"minimal_sns_app/model"

	"github.com/labstack/echo/v4"
)

// Validator

type UserValidatorMock struct {
	UserExistsResult bool
	Err              error
}

func (m *UserValidatorMock) UserExists(id string) (bool, error) {
	return m.UserExistsResult, m.Err
}

type ValidatorMock struct {
	Err error
}

func (v *ValidatorMock) ValidateID(id string) error {
	return v.Err
}

// UserRepository

type UserRepositoryMock struct {
	Exists    bool
	ExistsErr error
	CreateErr error
}

func (m *UserRepositoryMock) IsUserIDExists(id string) (bool, error) {
	return m.Exists, m.ExistsErr
}

func (m *UserRepositoryMock) CreateUser(id, name string) error {
	return m.CreateErr
}

// BlockRepository

type BlockRepositoryMock struct {
	IsBlockedResult  bool
	IsBlockedErr     error
	DeleteFriendErr  error
	RejectRequestErr error
	BlockErr         error
}

func (m *BlockRepositoryMock) IsBlocked(u1, u2 string) (bool, error) {
	return m.IsBlockedResult, m.IsBlockedErr
}

func (m *BlockRepositoryMock) DeleteFriendLink(u1, u2 string) error {
	return m.DeleteFriendErr
}

func (m *BlockRepositoryMock) RejectRequests(u1, u2 string) error {
	return m.RejectRequestErr
}

func (m *BlockRepositoryMock) Block(u1, u2 string) error {
	return m.BlockErr
}

// FriendRequestRepository

type FriendRequestRepositoryMock struct {
	IsBlockedResult       bool
	IsBlockedErr          error
	IsFriendResult        bool
	IsFriendErr           error
	HasReverseResult      bool
	HasReverseErr         error
	HasRequestedResult    bool
	HasRequestedErr       error
	RequestErr            error
	PendingRequestsResult []model.FriendRequest
	PendingRequestsErr    error
}

func (m *FriendRequestRepositoryMock) IsBlockedEachOther(u1, u2 string) (bool, error) {
	return m.IsBlockedResult, m.IsBlockedErr
}

func (m *FriendRequestRepositoryMock) IsAlreadyFriends(u1, u2 string) (bool, error) {
	return m.IsFriendResult, m.IsFriendErr
}

func (m *FriendRequestRepositoryMock) HasPendingRequest(u1, u2 string) (bool, error) {
	return m.HasReverseResult, m.HasReverseErr
}

func (m *FriendRequestRepositoryMock) HasAlreadyRequested(u1, u2 string) (bool, error) {
	return m.HasRequestedResult, m.HasRequestedErr
}

func (m *FriendRequestRepositoryMock) Request(u1, u2 string) error {
	return m.RequestErr
}

func (m *FriendRequestRepositoryMock) GetPendingRequests(userID string) ([]model.FriendRequest, error) {
	return m.PendingRequestsResult, m.PendingRequestsErr
}

// RespondFriend

type RespondRepositoryMock struct {
	FindRequestResult *model.FriendRequest
	FindRequestErr    error
	UpdateRequestErr  error
	CreateFriendErr   error
	RespondErr        error
}

func (m *RespondRepositoryMock) FindRequest(user1, user2 string) (*model.FriendRequest, error) {
	if m.FindRequestResult == nil {
		return &model.FriendRequest{}, m.FindRequestErr
	}
	return m.FindRequestResult, m.FindRequestErr
}

func (m *RespondRepositoryMock) UpdateRequest(req *model.FriendRequest, action string) error {
	return m.UpdateRequestErr
}

func (m *RespondRepositoryMock) CreateFriendLink(user1, user2 string) error {
	return m.CreateFriendErr
}

func (m *RespondRepositoryMock) RespondRequest(fromID, toID, action string) error {
	return m.RespondErr
}

// FriendRepository モック

type FriendRepositoryMock struct {
	Friends []model.Friend
	Err     error
}

func (m *FriendRepositoryMock) GetFriends(id string) ([]model.Friend, error) {
	return m.Friends, m.Err
}

type FriendOfFriendRepositoryMock struct {
	Result []model.Friend
	Err    error
}

func (m *FriendOfFriendRepositoryMock) GetFriendOfFriend(id string) ([]model.Friend, error) {
	return m.Result, m.Err
}

type FriendOfFriendPagingRepositoryMock struct {
	Result []model.Friend
	Err    error
}

func (m *FriendOfFriendPagingRepositoryMock) GetFriendOfFriendByIDWithPaging(id string, limit, offset int) ([]model.Friend, error) {
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
