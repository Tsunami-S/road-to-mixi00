package create

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"minimal_sns_app/model"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// --- モック定義 ---

type mockRespondRepo struct {
	findResult *model.FriendRequest
	findErr    error
	updateErr  error
	linkErr    error
}

func (m *mockRespondRepo) FindRequest(user1ID, user2ID string) (*model.FriendRequest, error) {
	return m.findResult, m.findErr
}
func (m *mockRespondRepo) UpdateRequest(req *model.FriendRequest, action string) error {
	return m.updateErr
}
func (m *mockRespondRepo) CreateFriendLink(user1ID, user2ID string) error {
	return m.linkErr
}
func (m *mockRespondRepo) RespondRequest(fromID, toID, action string) error {
	return nil
}

type mockValidator struct {
	Exists bool
	Err    error
}

func (v *mockValidator) UserExists(id string) (bool, error) {
	return v.Exists, v.Err
}

// --- 単体テスト ---

func TestRespondRequest(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name         string
		input        model.RespondRequestInput
		repo         *mockRespondRepo
		validator    *mockValidator
		wantCode     int
		wantContains string
	}{
		{
			name: "success: accepted",
			input: model.RespondRequestInput{
				User1ID: "user01", User2ID: "user02", Action: "accepted",
			},
			repo:      &mockRespondRepo{findResult: &model.FriendRequest{}},
			validator: &mockValidator{Exists: true},
			wantCode:  http.StatusOK, wantContains: "request accepted",
		},
		{
			name: "success: rejected",
			input: model.RespondRequestInput{
				User1ID: "user01", User2ID: "user02", Action: "rejected",
			},
			repo:      &mockRespondRepo{findResult: &model.FriendRequest{}},
			validator: &mockValidator{Exists: true},
			wantCode:  http.StatusOK, wantContains: "request rejected",
		},
		{
			name: "error: invalid user ID (empty)",
			input: model.RespondRequestInput{
				User1ID: "", User2ID: "user02", Action: "accepted",
			},
			repo:      &mockRespondRepo{},
			validator: &mockValidator{Exists: true},
			wantCode:  http.StatusBadRequest, wantContains: "invalid user IDs",
		},
		{
			name: "error: same user",
			input: model.RespondRequestInput{
				User1ID: "user01", User2ID: "user01", Action: "accepted",
			},
			repo:      &mockRespondRepo{},
			validator: &mockValidator{Exists: true},
			wantCode:  http.StatusBadRequest, wantContains: "invalid user IDs",
		},
		{
			name: "error: invalid action",
			input: model.RespondRequestInput{
				User1ID: "user01", User2ID: "user02", Action: "maybe",
			},
			repo:      &mockRespondRepo{},
			validator: &mockValidator{Exists: true},
			wantCode:  http.StatusBadRequest, wantContains: "invalid action",
		},
		{
			name: "error: user1 not found",
			input: model.RespondRequestInput{
				User1ID: "invalid", User2ID: "user02", Action: "accepted",
			},
			repo:      &mockRespondRepo{},
			validator: &mockValidator{Exists: false},
			wantCode:  http.StatusBadRequest, wantContains: "user1_id: user ID not found",
		},
		{
			name: "error: user2 not found",
			input: model.RespondRequestInput{
				User1ID: "user01", User2ID: "invalid", Action: "accepted",
			},
			repo:      &mockRespondRepo{},
			validator: &mockValidator{Exists: true},
			wantCode:  http.StatusBadRequest, wantContains: "user2_id: user ID not found",
		},
		{
			name: "error: request not found",
			input: model.RespondRequestInput{
				User1ID: "user01", User2ID: "user02", Action: "accepted",
			},
			repo:      &mockRespondRepo{findErr: errors.New("not found")},
			validator: &mockValidator{Exists: true},
			wantCode:  http.StatusBadRequest, wantContains: "request not found",
		},
		{
			name: "error: update failed",
			input: model.RespondRequestInput{
				User1ID: "user01", User2ID: "user02", Action: "accepted",
			},
			repo:      &mockRespondRepo{findResult: &model.FriendRequest{}, updateErr: errors.New("fail")},
			validator: &mockValidator{Exists: true},
			wantCode:  http.StatusInternalServerError, wantContains: "failed to update",
		},
		{
			name: "error: friend link creation failed",
			input: model.RespondRequestInput{
				User1ID: "user01", User2ID: "user02", Action: "accepted",
			},
			repo:      &mockRespondRepo{findResult: &model.FriendRequest{}, linkErr: errors.New("fail")},
			validator: &mockValidator{Exists: true},
			wantCode:  http.StatusInternalServerError, wantContains: "failed to create friendship",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.input)
			req := httptest.NewRequest(http.MethodPost, "/respond_friend_request", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := NewFriendRespondHandler(tc.validator, tc.repo)
			err := handler.RespondRequest(c)

			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantContains)
		})
	}
}
