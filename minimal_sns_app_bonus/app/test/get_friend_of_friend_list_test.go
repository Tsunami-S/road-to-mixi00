package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"minimal_sns_app/db"
	"minimal_sns_app/handler/get"
	"minimal_sns_app/handler/validate"
	repo_get "minimal_sns_app/repository/get"

	"github.com/labstack/echo/v4"
)

func setupTestDB_FOF_List(t *testing.T) {
	db.DB = initTestDB()
}

func TestGetFriendOfFriendList_Scenarios(t *testing.T) {
	setupTestDB_FOF_List(t)
	e := echo.New()

	handler := get.NewFriendOfFriendHandler(
		&validate.RealValidator{},
		&repo_get.RealFriendOfFriendRepository{},
	)

	tests := []struct {
		name      string
		userID    string
		wantCode  int
		wantBody  string
		notInBody string
	}{
		{
			name:     "1.友達の友達が返される（id1）",
			userID:   "id1",
			wantCode: http.StatusOK,
			wantBody: "user13",
		},
		{
			name:      "2.直接のフレンドは含まれない",
			userID:    "id1",
			wantCode:  http.StatusOK,
			notInBody: "user02",
		},
		{
			name:      "3.自分がブロックしている相手は含まれない",
			userID:    "id1",
			wantCode:  http.StatusOK,
			notInBody: "user18",
		},
		{
			name:      "4.自分がブロックされている相手は含まれない",
			userID:    "id1",
			wantCode:  http.StatusOK,
			notInBody: "user19",
		},
		{
			name:      "5.ブロックしている相手は含まれない",
			userID:    "id1",
			wantCode:  http.StatusOK,
			notInBody: "user16",
		},
		{
			name:      "6.ブロックされている相手は含まれない",
			userID:    "id1",
			wantCode:  http.StatusOK,
			notInBody: "user17",
		},
		{
			name:     "7.存在しないID",
			userID:   "invalid_id",
			wantCode: http.StatusBadRequest,
			wantBody: "user not found",
		},
		{
			name:     "8.友達もブロックもないユーザー",
			userID:   "id44",
			wantCode: http.StatusOK,
			wantBody: "no friends of friends found",
		},
		{
			name:     "9.相互関係の友達の友達は重複しない",
			userID:   "id1",
			wantCode: http.StatusOK,
			wantBody: "user11",
		},
		{
			name:     "10.共通の友人は重複しない",
			userID:   "id1",
			wantCode: http.StatusOK,
			wantBody: "user12",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/get_friend_of_friend_list?id="+tc.userID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if err := handler.FriendOfFriend(c); err != nil {
				t.Fatal(err)
			}

			body := rec.Body.String()

			if rec.Code != tc.wantCode {
				t.Errorf("ステータスコード不一致: got %d, want %d", rec.Code, tc.wantCode)
			}
			if tc.wantBody != "" && !strings.Contains(body, tc.wantBody) {
				t.Errorf("期待する文字列が含まれていない: want=%q, got=%q", tc.wantBody, body)
			}
			if tc.notInBody != "" && strings.Contains(body, tc.notInBody) {
				t.Errorf("含まれてはいけない文字列が含まれている: notWant=%q, got=%q", tc.notInBody, body)
			}
			if tc.name == "9.相互関係の友達の友達は重複しない" {
				count := strings.Count(body, "user14")
				if count > 1 {
					t.Errorf("user14 が重複して含まれている: 出現数 = %d", count)
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
