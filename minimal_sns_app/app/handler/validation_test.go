package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/db"
	"minimal_sns_app/test"
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
		{
			name:      "✅ 正常な値",
			limit:     "5",
			page:      "2",
			wantErr:   false,
			wantLimit: 5,
			wantPage:  2,
		},
		{
			name:    "❌ limit が負数",
			limit:   "-1",
			page:    "1",
			wantErr: true,
		},
		{
			name:    "❌ page がゼロ",
			limit:   "5",
			page:    "0",
			wantErr: true,
		},
		{
			name:    "❌ 数値でない",
			limit:   "abc",
			page:    "xyz",
			wantErr: true,
		},
		{
			name:    "❌ 空文字",
			limit:   "",
			page:    "",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/?limit="+tc.limit+"&page="+tc.page, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			limit, page, err := parseAndValidatePagination(c)
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
	db.DB = test.InitTestDB()

	tests := []struct {
		name       string
		userID     int
		wantExists bool
	}{
		{
			name:       "✅ 存在するユーザー",
			userID:     1,
			wantExists: true,
		},
		{
			name:       "❌ 存在しないユーザー",
			userID:     9999,
			wantExists: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			exists, err := userExists(tc.userID)
			if err != nil {
				t.Errorf("エラーが発生しました: %v", err)
			}
			if exists != tc.wantExists {
				t.Errorf("期待値と異なります: got=%v, want=%v", exists, tc.wantExists)
			}
		})
	}
}
