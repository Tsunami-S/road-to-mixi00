package test

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/db"
	"minimal_sns_app/handler"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestDB_Respond(t *testing.T) {
	db.DB = initTestDB()
}

func TestRespondFriendRequest_Scenarios(t *testing.T) {
	setupTestDB_Respond(t)
	e := echo.New()

	tests := []struct {
		name     string
		user1ID  string
		user2ID  string
		action   string
		wantCode int
		wantBody string
	}{
		{
			name:     "1.正常に申請を承認",
			user1ID:  "id45",
			user2ID:  "id1",
			action:   "accepted",
			wantCode: http.StatusOK,
			wantBody: "request accepted",
		},
		{
			name:     "2.正常に申請を拒否",
			user1ID:  "id11",
			user2ID:  "id27",
			action:   "rejected",
			wantCode: http.StatusOK,
			wantBody: "request rejected",
		},
		{
			name:     "3.自分自身に応答",
			user1ID:  "id1",
			user2ID:  "id1",
			action:   "accepted",
			wantCode: http.StatusBadRequest,
			wantBody: "invalid user IDs",
		},
		{
			name:     "4.存在しない申請",
			user1ID:  "id3",
			user2ID:  "id4",
			action:   "accepted",
			wantCode: http.StatusBadRequest,
			wantBody: "request not found",
		},
		{
			name:     "5.不正なアクション",
			user1ID:  "id45",
			user2ID:  "id1",
			action:   "maybe",
			wantCode: http.StatusBadRequest,
			wantBody: "invalid action",
		},
		{
			name:     "6.無効な user1_id",
			user1ID:  "invalid",
			user2ID:  "id1",
			action:   "accepted",
			wantCode: http.StatusBadRequest,
			wantBody: "user1_id: user ID not found",
		},
		{
			name:     "7.無効な user2_id",
			user1ID:  "id1",
			user2ID:  "invalid",
			action:   "accepted",
			wantCode: http.StatusBadRequest,
			wantBody: "user2_id: user ID not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url := "/respond_friend_request?user1_id=" + tc.user1ID +
				"&user2_id=" + tc.user2ID + "&action=" + tc.action

			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if err := handler.RespondFriendRequest(c); err != nil {
				t.Fatal(err)
			}

			if rec.Code != tc.wantCode {
				t.Errorf("ステータスコード不一致: got=%d, want=%d", rec.Code, tc.wantCode)
			}
			if !strings.Contains(rec.Body.String(), tc.wantBody) {
				t.Errorf("期待する文字列が含まれていない: want=%q, got=%q", tc.wantBody, rec.Body.String())
			}
		})
	}
}
