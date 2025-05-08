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

func TestRespondRequest(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name         string
		input        model.RespondRequestInput
		repo         *mockRespondRepo
		wantCode     int
		wantContains string
	}{
		{
			name: "success: accepted",
			input: model.RespondRequestInput{
				User1ID: "user01", User2ID: "user02", Action: "accepted",
			},
			repo:         &mockRespondRepo{findResult: &model.FriendRequest{}},
			wantCode:     http.StatusOK,
			wantContains: "request accepted",
		},
		{
			name: "success: rejected",
			input: model.RespondRequestInput{
				User1ID: "user01", User2ID: "user02", Action: "rejected",
			},
			repo:         &mockRespondRepo{findResult: &model.FriendRequest{}},
			wantCode:     http.StatusOK,
			wantContains: "request rejected",
		},
		{
			name: "error: invalid user ID",
			input: model.RespondRequestInput{
				User1ID: "", User2ID: "user02", Action: "accepted",
			},
			repo:         &mockRespondRepo{},
			wantCode:     http.StatusBadRequest,
			wantContains: "invalid user IDs",
		},
		{
			name: "error: same user",
			input: model.RespondRequestInput{
				User1ID: "user01", User2ID: "user01", Action: "accepted",
			},
			repo:         &mockRespondRepo{},
			wantCode:     http.StatusBadRequest,
			wantContains: "invalid user IDs",
		},
		{
			name: "error: invalid action",
			input: model.RespondRequestInput{
				User1ID: "user01", User2ID: "user02", Action: "maybe",
			},
			repo:         &mockRespondRepo{},
			wantCode:     http.StatusBadRequest,
			wantContains: "invalid action",
		},
		{
			name: "error: request not found",
			input: model.RespondRequestInput{
				User1ID: "user01", User2ID: "user02", Action: "accepted",
			},
			repo:         &mockRespondRepo{findErr: errors.New("not found")},
			wantCode:     http.StatusBadRequest,
			wantContains: "request not found",
		},
		{
			name: "error: update failed",
			input: model.RespondRequestInput{
				User1ID: "user01", User2ID: "user02", Action: "accepted",
			},
			repo:         &mockRespondRepo{findResult: &model.FriendRequest{}, updateErr: errors.New("update fail")},
			wantCode:     http.StatusInternalServerError,
			wantContains: "failed to update",
		},
		{
			name: "error: link creation failed",
			input: model.RespondRequestInput{
				User1ID: "user01", User2ID: "user02", Action: "accepted",
			},
			repo:         &mockRespondRepo{findResult: &model.FriendRequest{}, linkErr: errors.New("link fail")},
			wantCode:     http.StatusInternalServerError,
			wantContains: "failed to create friendship",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.input)
			req := httptest.NewRequest(http.MethodPost, "/respond_friend_request", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := NewFriendRespondHandler(nil, tc.repo)
			err := handler.RespondRequest(c)

			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantContains)
		})
	}
}
