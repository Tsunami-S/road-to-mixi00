package get

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"minimal_sns_app/model"
	"minimal_sns_app/test/mock"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPendingRequestHandler(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name         string
		userID       string
		valExists    bool
		valErr       error
		repoResult   []model.FriendRequest
		repoErr      error
		wantCode     int
		wantContains string
	}{
		{
			name:         "正常系: 保留中リクエストあり",
			userID:       "user01",
			valExists:    true,
			repoResult:   []model.FriendRequest{{User1ID: "user01", User2ID: "user02", Status: "pending"}},
			wantCode:     http.StatusOK,
			wantContains: "user02",
		},
		{
			name:         "正常系: リクエストなし",
			userID:       "user01",
			valExists:    true,
			repoResult:   []model.FriendRequest{},
			wantCode:     http.StatusOK,
			wantContains: "no pending requests found",
		},
		{
			name:         "異常系: 存在しないユーザー",
			userID:       "userX",
			valExists:    false,
			wantCode:     http.StatusBadRequest,
			wantContains: "user not found",
		},
		{
			name:         "異常系: Validatorエラー",
			userID:       "user01",
			valErr:       errors.New("validation error"),
			wantCode:     http.StatusBadRequest,
			wantContains: "validation error",
		},
		{
			name:         "異常系: Repositoryエラー",
			userID:       "user01",
			valExists:    true,
			repoErr:      errors.New("db error"),
			wantCode:     http.StatusInternalServerError,
			wantContains: "db error",
		},
		{
			name:         "異常系: user_id 空",
			userID:       "",
			wantCode:     http.StatusBadRequest,
			wantContains: "user_id must be",
		},
		{
			name:         "異常系: user_id 長すぎ",
			userID:       strings.Repeat("a", 21),
			wantCode:     http.StatusBadRequest,
			wantContains: "user_id must be",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := &mock.UserValidatorMock{
				UserExistsResult: tc.valExists,
				Err:              tc.valErr,
			}
			r := &mock.FriendRequestRepositoryMock{
				PendingRequestsResult: tc.repoResult,
				PendingRequestsErr:    tc.repoErr,
			}
			h := NewPendingRequestHandler(v, r)

			req := httptest.NewRequest(http.MethodGet, "/?user_id="+tc.userID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.PendingRequests(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantContains)
		})
	}
}
