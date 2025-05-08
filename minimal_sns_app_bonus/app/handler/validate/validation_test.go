package validate

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRealValidator_UserExists(t *testing.T) {
	tests := []struct {
		name       string
		inputID    string
		mockFunc   func(string) (bool, error)
		wantValid  bool
		wantErrStr string
	}{
		{
			name:    "正常系: ユーザー存在",
			inputID: "user01",
			mockFunc: func(id string) (bool, error) {
				return true, nil
			},
			wantValid: true,
		},
		{
			name:       "異常系: 空ID",
			inputID:    "",
			wantValid:  false,
			wantErrStr: "invalid user ID format",
		},
		{
			name:       "異常系: 長すぎるID",
			inputID:    "123456789012345678901",
			wantValid:  false,
			wantErrStr: "invalid user ID format",
		},
		{
			name:    "異常系: DBエラー",
			inputID: "user01",
			mockFunc: func(id string) (bool, error) {
				return false, errors.New("some DB error")
			},
			wantValid:  false,
			wantErrStr: "DB error",
		},
		{
			name:    "異常系: ユーザーが存在しない",
			inputID: "user01",
			mockFunc: func(id string) (bool, error) {
				return false, nil
			},
			wantValid:  false,
			wantErrStr: "user not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			validator := &RealValidator{
				UserExistsFunc: tc.mockFunc,
			}

			ok, err := validator.UserExists(tc.inputID)
			assert.Equal(t, tc.wantValid, ok)

			if tc.wantErrStr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErrStr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRealPaginationValidator_ParseAndValidatePagination(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name       string
		query      string
		wantLimit  int
		wantPage   int
		wantStatus int
		wantErrMsg string
	}{
		{
			name:      "正常系: limit=10, page=2",
			query:     "?limit=10&page=2",
			wantLimit: 10,
			wantPage:  2,
		},
		{
			name:       "異常系: limitが0",
			query:      "?limit=0&page=2",
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "invalid limit",
		},
		{
			name:       "異常系: pageがマイナス",
			query:      "?limit=10&page=-1",
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "invalid page",
		},
		{
			name:       "異常系: limitが文字列",
			query:      "?limit=abc&page=2",
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "invalid limit",
		},
		{
			name:       "異常系: pageが空",
			query:      "?limit=10&page=",
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "invalid page",
		},
		{
			name:       "異常系: limitが空",
			query:      "?limit=&page=1",
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "invalid limit",
		},
	}

	validator := &RealPaginationValidator{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test"+tc.query, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			limit, page, err := validator.ParseAndValidatePagination(c)

			if tc.wantErrMsg != "" {
				assert.Error(t, err)
				he, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tc.wantStatus, he.Code)
				assert.Contains(t, he.Message.(string), tc.wantErrMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantLimit, limit)
				assert.Equal(t, tc.wantPage, page)
			}
		})
	}
}
