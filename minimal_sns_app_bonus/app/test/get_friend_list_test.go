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

func setupTestDB_FriendList(t *testing.T) {
	db.DB = initTestDB()
}

func TestGetFriendList_Scenarios(t *testing.T) {
	setupTestDB_FriendList(t)
	e := echo.New()

	handler := get.NewFriendHandler(&validate.RealValidator{}, &repo_get.RealFriendRepository{})

	tests := []struct {
		name      string
		userID    string
		wantCode  int
		wantBody  string
		notInBody string
	}{
		{
			name:     "1.自分からのフレンドリンク",
			userID:   "id1",
			wantCode: http.StatusOK,
			wantBody: "user02",
		},
		{
			name:     "2.相手からのフレンドリンク",
			userID:   "id1",
			wantCode: http.StatusOK,
			wantBody: "user04",
		},
		{
			name:      "3.ブロックしているユーザーは含めない",
			userID:    "id1",
			wantCode:  http.StatusOK,
			notInBody: "user39",
		},
		{
			name:      "4.相手からブロックされているユーザーは含めない",
			userID:    "id1",
			wantCode:  http.StatusOK,
			notInBody: "user40",
		},
		{
			name:      "5.無関係のユーザーは含めない",
			userID:    "id1",
			wantCode:  http.StatusOK,
			notInBody: "user44",
		},
		{
			name:     "6.存在しないID",
			userID:   "invalid_id",
			wantCode: http.StatusBadRequest,
			wantBody: "user not found",
		},
		{
			name:      "7.一方的にブロックされているユーザーは除外される",
			userID:    "id6",
			wantCode:  http.StatusOK,
			notInBody: "user03",
		},
		{
			name:     "8.フレンドもブロックもない新規ユーザー",
			userID:   "id44",
			wantCode: http.StatusOK,
			wantBody: "no friends found",
		},
		{
			name:     "9.相互にフレンド登録されたユーザーは重複しない",
			userID:   "id1",
			wantCode: http.StatusOK,
			wantBody: "user10",
		},
		{
			name:      "10.自分自身へのフレンド",
			userID:    "id1",
			wantCode:  http.StatusOK,
			notInBody: "user01",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/get_friend_list?id="+tc.userID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if err := handler.Friend(c); err != nil {
				t.Fatal(err)
			}

			body := rec.Body.String()

			if rec.Code != tc.wantCode {
				t.Errorf("ステータスコード不一致: got=%d, want=%d", rec.Code, tc.wantCode)
			}
			if tc.wantBody != "" && !strings.Contains(body, tc.wantBody) {
				t.Errorf("期待する文字列が含まれていない: want=%q, got=%q", tc.wantBody, body)
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
