package get

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"minimal_sns_app/model"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockValidator struct {
	exists bool
	err    error
}

func (m *mockValidator) UserExists(id string) (bool, error) {
	return m.exists, m.err
}

type mockFriendRepository struct {
	friends []model.Friend
	err     error
}

func (m *mockFriendRepository) GetFriends(id string) ([]model.Friend, error) {
	return m.friends, m.err
}

func TestFriendHandler(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name          string
		userID        string
		mockExists    bool
		mockExistErr  error
		mockFriends   []model.Friend
		mockFriendErr error
		wantCode      int
		wantBodyPart  string
	}{
		{
			name:         "正常系：フレンドあり",
			userID:       "user01",
			mockExists:   true,
			mockFriends:  []model.Friend{{ID: "u2", Name: "user02"}},
			wantCode:     http.StatusOK,
			wantBodyPart: "user02",
		},
		{
			name:         "正常系：フレンドなし",
			userID:       "user01",
			mockExists:   true,
			mockFriends:  []model.Friend{},
			wantCode:     http.StatusOK,
			wantBodyPart: "no friends found",
		},
		{
			name:         "異常系：存在しないユーザー",
			userID:       "userX",
			mockExists:   false,
			wantCode:     http.StatusBadRequest,
			wantBodyPart: "user not found",
		},
		{
			name:         "異常系：validatorでエラー",
			userID:       "userX",
			mockExistErr: errors.New("DB error"),
			wantCode:     http.StatusInternalServerError,
			wantBodyPart: "DB error",
		},
		{
			name:          "異常系：friend取得でエラー",
			userID:        "userX",
			mockExists:    true,
			mockFriendErr: errors.New("repo error"),
			wantCode:      http.StatusInternalServerError,
			wantBodyPart:  "repo error",
		},
		{
			name:         "異常系：userIDが空",
			userID:       "",
			wantCode:     http.StatusBadRequest,
			wantBodyPart: "id must be",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := &mockValidator{exists: tc.mockExists, err: tc.mockExistErr}
			r := &mockFriendRepository{friends: tc.mockFriends, err: tc.mockFriendErr}
			h := NewFriendHandler(v, r)

			req := httptest.NewRequest(http.MethodGet, "/?id="+tc.userID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.Friend(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantBodyPart)
		})
	}
}
