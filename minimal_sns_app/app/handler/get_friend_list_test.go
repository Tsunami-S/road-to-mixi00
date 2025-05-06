package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/db"
	"minimal_sns_app/test"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestDB(t *testing.T) {
	db.DB = test.InitTestDB()
}

func TestGetFriendList_Scenarios(t *testing.T) {
	setupTestDB(t)
	e := echo.New()

	tests := []struct {
		name      string
		userID    string
		wantCode  int
		wantBody  string
		notInBody string
	}{
		{
			name:     "✅ 自分からのフレンドリンク",
			userID:   "1",
			wantCode: http.StatusOK,
			wantBody: "user02",
		},
		{
			name:     "✅ 相手からのフレンドリンク",
			userID:   "1",
			wantCode: http.StatusOK,
			wantBody: "user04",
		},
		{
			name:      "❌ ブロックしているユーザーは含めない",
			userID:    "1",
			wantCode:  http.StatusOK,
			notInBody: "user39",
		},
		{
			name:      "❌ 相手からブロックされているユーザーは含めない",
			userID:    "1",
			wantCode:  http.StatusOK,
			notInBody: "user40",
		},
		{
			name:      "❌ 無関係のユーザーは含めない",
			userID:    "1",
			wantCode:  http.StatusOK,
			notInBody: "user44",
		},
		{
			name:     "❌ 存在しないID",
			userID:   "9999",
			wantCode: http.StatusBadRequest,
			wantBody: "user not found",
		},
		{
			name:      "🟥 一方的にブロックされているユーザーは除外される",
			userID:    "6",
			wantCode:  http.StatusOK,
			notInBody: "user03",
		},
		{
			name:     "🟩 フレンドもブロックもない新規ユーザー",
			userID:   "44",
			wantCode: http.StatusOK,
			wantBody: "no friends found",
		},
		{
			name:     "🟦 相互にフレンド登録されたユーザーは重複しない",
			userID:   "1",
			wantCode: http.StatusOK,
			wantBody: "user10", // 後で count で重複確認
		},
		{
			name:      "❌ 自分自身へのフレンド",
			userID:    "1",
			wantCode:  http.StatusOK,
			notInBody: "user01",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/get_friend_list?id="+tc.userID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if err := GetFriendList(c); err != nil {
				t.Fatal(err)
			}

			if rec.Code != tc.wantCode {
				t.Errorf("ステータスコード不一致: got %d, want %d", rec.Code, tc.wantCode)
			}

			body := rec.Body.String()
			if tc.wantBody != "" && !strings.Contains(body, tc.wantBody) {
				t.Errorf("期待する文字列が含まれない: want=%q, got=%q", tc.wantBody, body)
			}
			if tc.notInBody != "" && strings.Contains(body, tc.notInBody) {
				t.Errorf("含まれてはいけない文字列が含まれている: notWant=%q, got=%q", tc.notInBody, body)
			}

			// 特別処理: 重複チェック（相互フレンドの重複防止）
			if tc.name == "🟦 相互にフレンド登録されたユーザーは重複しない" {
				count := strings.Count(body, "user10")
				if count > 1 {
					t.Errorf("user10 が重複して含まれている: 出現数 = %d\nレスポンス: %s", count, body)
				}
			}
		})
	}
}
