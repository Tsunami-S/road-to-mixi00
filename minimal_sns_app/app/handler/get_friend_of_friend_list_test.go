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

func setupTestDB_FOF(t *testing.T) {
	db.DB = test.InitTestDB()
}

func TestGetFriendOfFriendList_Scenarios(t *testing.T) {
	setupTestDB_FOF(t)
	e := echo.New()

	tests := []struct {
		name      string
		userID    string
		wantCode  int
		wantBody  string
		notInBody string
	}{
		{
			name:     "✅ 友達の友達が返される（id1）",
			userID:   "1",
			wantCode: http.StatusOK,
			wantBody: "user13", // id2 のフレンド → 友達の友達
		},
		{
			name:      "❌ 直接のフレンドは含まれない",
			userID:    "1",
			wantCode:  http.StatusOK,
			notInBody: "user02", // id2 は直接のフレンド
		},
		{
			name:      "❌ 自分がブロックしている相手は含まれない",
			userID:    "1",
			wantCode:  http.StatusOK,
			notInBody: "user18",
		},
		{
			name:      "❌ 自分がブロックされている相手は含まれない",
			userID:    "1",
			wantCode:  http.StatusOK,
			notInBody: "user19",
		},
		{
			name:      "❌ ブロックしている相手は含まれない",
			userID:    "1",
			wantCode:  http.StatusOK,
			notInBody: "user16",
		},
		{
			name:      "❌ ブロックされている相手は含まれない",
			userID:    "1",
			wantCode:  http.StatusOK,
			notInBody: "user17",
		},
		{
			name:     "❌ 存在しないID",
			userID:   "9999",
			wantCode: http.StatusBadRequest,
			wantBody: "user not found",
		},
		{
			name:     "🟩 友達もブロックもないユーザー",
			userID:   "44",
			wantCode: http.StatusOK,
			wantBody: "no friends of friends found",
		},
		{
			name:     "🟦 相互関係の友達の友達は重複しない",
			userID:   "1",
			wantCode: http.StatusOK,
			wantBody: "user11",
		},
		{
			name:     "🟨 共通の友人は重複しない",
			userID:   "1",
			wantCode: http.StatusOK,
			wantBody: "user12",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/get_friend_of_friend_list?id="+tc.userID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if err := GetFriendOfFriendList(c); err != nil {
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
			if tc.name == "🟦 相互関係の友達の友達は重複しない" {
				count := strings.Count(rec.Body.String(), "user14")
				if count > 1 {
					t.Errorf("user14 が重複して含まれている: 出現数 = %d", count)
				}
			}
			if tc.name == "🟨 共通の友人は重複しない" {
				count := strings.Count(rec.Body.String(), "user12")
				if count > 1 {
					t.Errorf("user12 が重複して含まれている: 出現数 = %d", count)
				}
			}
		})
	}
}
