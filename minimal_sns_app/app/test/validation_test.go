package test

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/db"
	"minimal_sns_app/handler"
	"minimal_sns_app/repository"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseAndValidatePagination(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name      string
		limit     string
		page      string
		wantErr   bool
		wantLimit int
		wantPage  int
	}{
		{"1.正常な値", "5", "2", false, 5, 2},
		{"2.limit が負数", "-1", "1", true, 0, 0},
		{"3.page がゼロ", "5", "0", true, 0, 0},
		{"4.数値でない", "abc", "xyz", true, 0, 0},
		{"5.空文字", "", "", true, 0, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/?limit="+tc.limit+"&page="+tc.page, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			limit, page, err := handler.ParseAndValidatePagination(c) // ← handler経由
			if tc.wantErr && err == nil {
				t.Errorf("期待したエラーが返らなかった")
			}
			if !tc.wantErr {
				if err != nil {
					t.Errorf("予期しないエラー: %v", err)
				}
				if limit != tc.wantLimit || page != tc.wantPage {
					t.Errorf("値が一致しない: got limit=%d page=%d, want limit=%d page=%d",
						limit, page, tc.wantLimit, tc.wantPage)
				}
			}
		})
	}
}

func TestUserExists(t *testing.T) {
	db.DB = InitTestDB()

	tests := []struct {
		name       string
		userID     int
		wantExists bool
	}{
		{
			name:       "1.存在するユーザー（id=1）",
			userID:     1,
			wantExists: true,
		},
		{
			name:       "2.存在しないユーザー（id=9999）",
			userID:     9999,
			wantExists: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			exists, err := repository.UserExists(tc.userID)
			if err != nil {
				t.Errorf("エラーが発生しました: %v", err)
			}
			if exists != tc.wantExists {
				t.Errorf("存在チェック結果が期待と異なります: got=%v, want=%v", exists, tc.wantExists)
			}
		})
	}
}
