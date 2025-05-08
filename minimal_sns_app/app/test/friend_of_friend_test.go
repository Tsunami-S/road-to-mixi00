package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"minimal_sns_app/handler/get"
	repo_get "minimal_sns_app/repository/get"
	"minimal_sns_app/repository/validate"

	"github.com/labstack/echo/v4"
)

func TestFriendHandler_Friend(t *testing.T) {
	setupTestDB(t)

	e := echo.New()

	handler := get.NewFriendHandler(
		&validate.RealValidator{},
		&repo_get.RealFriendRepository{},
	)

	tests := []struct {
		name      string
		queryID   string
		wantCode  int
		wantBody  string
		notInBody string
		checkDup  string
	}{
		{"フレンドリンク（自分発）", "1", http.StatusOK, "user02", "", ""},
		{"フレンドリンク（相手発）", "1", http.StatusOK, "user04", "", ""},
		{"ブロックしたユーザーを除外", "1", http.StatusOK, "", "user39", ""},
		{"ブロックされたユーザーを除外", "1", http.StatusOK, "", "user40", ""},
		{"無関係ユーザーを除外", "1", http.StatusOK, "", "user44", ""},
		{"存在しないID", "9999", http.StatusBadRequest, "user not found", "", ""},
		{"ブロックされている場合除外", "6", http.StatusOK, "", "user03", ""},
		{"新規ユーザー（フレンド・ブロックなし）", "44", http.StatusOK, "no friends found", "", ""},
		{"相互フレンドは重複しない", "1", http.StatusOK, "user10", "", "user10"},
		{"自分自身を除外", "1", http.StatusOK, "", "user01", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/get_friend_list?id="+tc.queryID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.Friend(c)
			if err != nil {
				t.Fatalf("ハンドラー内でエラー発生: %v", err)
			}

			body := rec.Body.String()

			if rec.Code != tc.wantCode {
				t.Errorf("ステータスコード不一致: got=%d, want=%d", rec.Code, tc.wantCode)
			}

			if tc.wantBody != "" && !strings.Contains(body, tc.wantBody) {
				t.Errorf("期待する文字列が含まれない: want=%q, got=%q", tc.wantBody, body)
			}

			if tc.notInBody != "" && strings.Contains(body, tc.notInBody) {
				t.Errorf("含んではいけない文字列が含まれている: notWant=%q, got=%q", tc.notInBody, body)
			}

			if tc.checkDup != "" {
				count := strings.Count(body, tc.checkDup)
				if count > 1 {
					t.Errorf("%s が重複して含まれている: 出現数=%d\nレスポンス: %s", tc.checkDup, count, body)
				}
			}
		})
	}
}
