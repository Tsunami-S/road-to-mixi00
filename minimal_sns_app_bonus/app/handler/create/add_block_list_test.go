package create

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"minimal_sns_app/model"
	"minimal_sns_app/test/mock"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAddBlockList(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name           string
		input          model.BlockRequest
		validatorErr   error
		validatorExist bool
		repo           *mock.BlockRepositoryMock
		wantCode       int
		wantContains   string
	}{
		{
			name:           "正常系",
			input:          model.BlockRequest{User1ID: "user01", User2ID: "user02"},
			validatorExist: true,
			repo:           &mock.BlockRepositoryMock{},
			wantCode:       http.StatusOK,
			wantContains:   "user blocked",
		},
		{
			name:         "異常系: 自分自身をブロック",
			input:        model.BlockRequest{User1ID: "user01", User2ID: "user01"},
			wantCode:     http.StatusBadRequest,
			wantContains: "cannot block yourself",
		},
		{
			name:         "異常系: user_id 長すぎ",
			input:        model.BlockRequest{User1ID: "u123456789012345678901", User2ID: "user02"},
			wantCode:     http.StatusBadRequest,
			wantContains: "invalid user IDs",
		},
		{
			name:         "異常系: user_id 空",
			input:        model.BlockRequest{User1ID: "", User2ID: "user02"},
			wantCode:     http.StatusBadRequest,
			wantContains: "invalid user IDs",
		},
		{
			name:           "異常系: user1 が存在しない",
			input:          model.BlockRequest{User1ID: "x", User2ID: "y"},
			validatorExist: false,
			wantCode:       http.StatusBadRequest,
			wantContains:   "user not found",
		},
		{
			name:         "異常系: validator error",
			input:        model.BlockRequest{User1ID: "x", User2ID: "y"},
			validatorErr: errors.New("user not found"),
			wantCode:     http.StatusBadRequest,
			wantContains: "user not found",
		},
		{
			name:           "異常系: 既にブロック済み",
			input:          model.BlockRequest{User1ID: "x", User2ID: "y"},
			validatorExist: true,
			repo:           &mock.BlockRepositoryMock{IsBlockedResult: true},
			wantCode:       http.StatusBadRequest,
			wantContains:   "already blocked",
		},
		{
			name:           "異常系: friendship 削除失敗",
			input:          model.BlockRequest{User1ID: "x", User2ID: "y"},
			validatorExist: true,
			repo:           &mock.BlockRepositoryMock{DeleteFriendErr: errors.New("delete error")},
			wantCode:       http.StatusInternalServerError,
			wantContains:   "failed to delete friendship",
		},
		{
			name:           "異常系: friend request 拒否失敗",
			input:          model.BlockRequest{User1ID: "x", User2ID: "y"},
			validatorExist: true,
			repo:           &mock.BlockRepositoryMock{RejectRequestErr: errors.New("reject error")},
			wantCode:       http.StatusInternalServerError,
			wantContains:   "failed to reject friend request",
		},
		{
			name:           "異常系: block 失敗",
			input:          model.BlockRequest{User1ID: "x", User2ID: "y"},
			validatorExist: true,
			repo:           &mock.BlockRepositoryMock{BlockErr: errors.New("block error")},
			wantCode:       http.StatusInternalServerError,
			wantContains:   "failed to block user",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			payload, _ := json.Marshal(tc.input)
			req := httptest.NewRequest(http.MethodPost, "/block", bytes.NewReader(payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			validator := &mock.UserValidatorMock{
				UserExistsResult: tc.validatorExist,
				Err:              tc.validatorErr,
			}
			repo := tc.repo
			if repo == nil {
				repo = &mock.BlockRepositoryMock{}
			}
			h := NewBlockHandler(validator, repo)
			err := h.AddBlockList(c)

			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantContains)
		})
	}
}
