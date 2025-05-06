package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/labstack/echo/v4"
	"minimal_sns_app/handler"
)

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
		{"1.自分からのフレンドリンク", "1", http.StatusOK, "user02", ""},
		{"2.相手からのフレンドリンク", "1", http.StatusOK, "user04", ""},
		{"3.ブロックしているユーザーは含めない", "1", http.StatusOK, "", "user39"},
		{"4.相手からブロックされているユーザーは含めない", "1", http.StatusOK, "", "user40"},
		{"5.無関係のユーザーは含めない", "1", http.StatusOK, "", "user44"},
		{"6.存在しないID", "9999", http.StatusBadRequest, "user not found", ""},
		{"7.一方的にブロックされているユーザーは除外される", "6", http.StatusOK, "", "user03"},
		{"8.フレンドもブロックもない新規ユーザー", "44", http.StatusOK, "no friends found", ""},
		{"9.相互にフレンド登録されたユーザーは重複しない", "1", http.StatusOK, "user10", ""},
		{"10.自分自身へのフレンド", "1", http.StatusOK, "", "user01"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/get_friend_list?id="+tc.userID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if err := handler.GetFriendList(c); err != nil {
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

			if tc.name == "9.相互にフレンド登録されたユーザーは重複しない" {
				count := strings.Count(body, "user10")
				if count > 1 {
					t.Errorf("user10 が重複して含まれている: 出現数 = %d\nレスポンス: %s", count, body)
				}
			}
		})
	}
}
