package get

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"minimal_sns_app/model"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockFriendOfFriendRepo struct {
	result []model.Friend
	err    error
}

func (m *mockFriendOfFriendRepo) GetFriendOfFriend(id string) ([]model.Friend, error) {
	return m.result, m.err
}

func TestFriendOfFriendHandler(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name         string
		userID       string
		mockExists   bool
		mockValErr   error
		mockResult   []model.Friend
		mockRepoErr  error
		wantCode     int
		wantContains string
	}{
		{
			name:         "正常系: フレンドのフレンドあり",
			userID:       "user01",
			mockExists:   true,
			mockResult:   []model.Friend{{ID: "u2", Name: "user02"}},
			wantCode:     http.StatusOK,
			wantContains: "user02",
		},
		{
			name:         "正常系: フレンドのフレンドなし",
			userID:       "user01",
			mockExists:   true,
			mockResult:   []model.Friend{},
			wantCode:     http.StatusOK,
			wantContains: "no friends of friends found",
		},
		{
			name:         "異常系: ユーザー存在しない",
			userID:       "not_found",
			mockExists:   false,
			wantCode:     http.StatusBadRequest,
			wantContains: "user not found",
		},
		{
			name:         "異常系: Validatorエラー",
			userID:       "userX",
			mockValErr:   errors.New("validation failed"),
			wantCode:     http.StatusInternalServerError,
			wantContains: "validation failed",
		},
		{
			name:         "異常系: Repoエラー",
			userID:       "user01",
			mockExists:   true,
			mockRepoErr:  errors.New("repo error"),
			wantCode:     http.StatusInternalServerError,
			wantContains: "repo error",
		},
		{
			name:         "異常系: 空ID",
			userID:       "",
			wantCode:     http.StatusBadRequest,
			wantContains: "user_id must be",
		},
		{
			name:         "異常系: ID長すぎ",
			userID:       strings.Repeat("a", 21),
			wantCode:     http.StatusBadRequest,
			wantContains: "user_id must be",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := &mockValidator{exists: tc.mockExists, err: tc.mockValErr}
			r := &mockFriendOfFriendRepo{result: tc.mockResult, err: tc.mockRepoErr}
			h := NewFriendOfFriendHandler(v, r)

			req := httptest.NewRequest(http.MethodGet, "/?id="+tc.userID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.FriendOfFriend(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantContains)
		})
	}
}
