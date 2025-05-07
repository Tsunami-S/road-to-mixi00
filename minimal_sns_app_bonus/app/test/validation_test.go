package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"minimal_sns_app/db"
	handle_valid "minimal_sns_app/handler/validate"
	repo_valid "minimal_sns_app/repository/validate"

	"github.com/labstack/echo/v4"
)

func setupTestDB_Validation(t *testing.T) {
	db.DB = initTestDB()
}

func TestIsValidUserId(t *testing.T) {
	setupTestDB_Validation(t)

	tests := []struct {
		name    string
		userID  string
		wantOK  bool
		wantErr string
	}{
		{
			name:   "1.有効なID",
			userID: "id1",
			wantOK: true,
		},
		{
			name:    "2.空文字",
			userID:  "",
			wantOK:  false,
			wantErr: "invalid user ID format",
		},
		{
			name:    "3.長すぎるID",
			userID:  "1234567890123456789012345",
			wantOK:  false,
			wantErr: "invalid user ID format",
		},
		{
			name:    "4.存在しないID",
			userID:  "invalid_id",
			wantOK:  false,
			wantErr: "user ID not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := handle_valid.IsValidUserId(tc.userID)
			if ok != tc.wantOK {
				t.Errorf("期待結果: %v, 実際: %v", tc.wantOK, ok)
			}
			if err != nil && tc.wantErr != "" && err.Error() != tc.wantErr {
				t.Errorf("エラーメッセージ不一致: got=%q, want=%q", err.Error(), tc.wantErr)
			}
		})
	}
}

func TestUserExists(t *testing.T) {
	setupTestDB_Validation(t)

	tests := []struct {
		name     string
		userID   string
		wantBool bool
	}{
		{
			name:     "1.存在するID",
			userID:   "id1",
			wantBool: true,
		},
		{
			name:     "2.存在しないID",
			userID:   "notfound",
			wantBool: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := repo_valid.UserExists(tc.userID)
			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
			if ok != tc.wantBool {
				t.Errorf("期待: %v, 実際: %v", tc.wantBool, ok)
			}
		})
	}
}

func TestParseAndValidatePagination(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name       string
		limit      string
		page       string
		wantLimit  int
		wantPage   int
		shouldFail bool
	}{
		{
			name:      "1.正常な入力",
			limit:     "5",
			page:      "2",
			wantLimit: 5,
			wantPage:  2,
		},
		{
			name:       "2.limitが負数",
			limit:      "-1",
			page:       "1",
			shouldFail: true,
		},
		{
			name:       "3.pageがゼロ",
			limit:      "5",
			page:       "0",
			shouldFail: true,
		},
		{
			name:       "4.不正な文字列",
			limit:      "abc",
			page:       "xyz",
			shouldFail: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/?limit="+tc.limit+"&page="+tc.page, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			limit, page, err := handle_valid.ParseAndValidatePagination(c)
			if tc.shouldFail {
				if err == nil {
					t.Errorf("エラーが期待されたが nil が返された")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}
				if limit != tc.wantLimit || page != tc.wantPage {
					t.Errorf("limit/page が一致しない: got (%d, %d), want (%d, %d)",
						limit, page, tc.wantLimit, tc.wantPage)
				}
			}
		})
	}
}
