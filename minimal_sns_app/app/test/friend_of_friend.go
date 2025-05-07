package test

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/handler/get"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
		{"1.友達の友達が返される（id1）", "1", http.StatusOK, "user13", ""},
		{"2.直接のフレンドは含まれない", "1", http.StatusOK, "", "user02"},
		{"3.自分がブロックしている相手は含まれない", "1", http.StatusOK, "", "user18"},
		{"4.自分がブロックされている相手は含まれない", "1", http.StatusOK, "", "user19"},
		{"5.ブロックしている相手は含まれない", "1", http.StatusOK, "", "user16"},
		{"6.ブロックされている相手は含まれない", "1", http.StatusOK, "", "user17"},
		{"7.存在しないID", "9999", http.StatusBadRequest, "user not found", ""},
		{"8.友達もブロックもないユーザー", "44", http.StatusOK, "no friends of friends found", ""},
		{"9.相互関係の友達の友達は重複しない", "1", http.StatusOK, "user11", ""},
		{"10.共通の友人は重複しない", "1", http.StatusOK, "user12", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/get_friend_of_friend_list?id="+tc.userID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if err := get.FriendOfFriend(c); err != nil {
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

			if tc.name == "9.相互関係の友達の友達は重複しない" {
				count := strings.Count(body, "user11")
				if count > 1 {
					t.Errorf("user11 が重複して含まれている: 出現数 = %d", count)
				}
			}
			if tc.name == "10.共通の友人は重複しない" {
				count := strings.Count(body, "user12")
				if count > 1 {
					t.Errorf("user12 が重複して含まれている: 出現数 = %d", count)
				}
			}
		})
	}
}

