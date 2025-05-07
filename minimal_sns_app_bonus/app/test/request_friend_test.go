package test

import (
	"minimal_sns_app/db"
	"minimal_sns_app/handler/create"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func setupTestDB_Request(t *testing.T) {
	db.DB = initTestDB()
}

func TestRequestFriend_Scenarios(t *testing.T) {
	setupTestDB_Request(t)
	e := echo.New()

	tests := []struct {
		name     string
		body     string
		wantCode int
		wantBody string
	}{
		{
			name:     "1.正常なフレンド申請",
			body:     `{"user1_id":"id32","user2_id":"id43"}`,
			wantCode: http.StatusOK,
			wantBody: "friend request sent",
		},
		{
			name:     "2.自分自身に申請",
			body:     `{"user1_id":"id1","user2_id":"id1"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "cannot request yourself",
		},
		{
			name:     "3.ブロックしている",
			body:     `{"user1_id":"id1","user2_id":"id39"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "cannot send friend request due to block",
		},
		{
			name:     "4.ブロックされている",
			body:     `{"user1_id":"id1","user2_id":"id40"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "cannot send friend request due to block",
		},
		{
			name:     "5.すでにフレンド",
			body:     `{"user1_id":"id1","user2_id":"id2"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "you are already friends",
		},
		{
			name:     "6.すでにフレンド(相手から)",
			body:     `{"user1_id":"id1","user2_id":"id4"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "you are already friends",
		},
		{
			name:     "7.逆方向にpendingな申請がある",
			body:     `{"user1_id":"id41","user2_id":"id1"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "you already have a pending request from this user",
		},
		{
			name:     "8.同じ方向の申請がすでにある",
			body:     `{"user1_id":"id1","user2_id":"id27"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "friend request already sent",
		},
		{
			name:     "9.存在しない user1_id",
			body:     `{"user1_id":"invalid_user","user2_id":"id2"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "user1_id: user ID not found",
		},
		{
			name:     "10.存在しない user2_id",
			body:     `{"user1_id":"id2","user2_id":"invalid_user"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "user2_id: user ID not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/request_friend", strings.NewReader(tc.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if err := create.RequestFriend(c); err != nil {
				t.Fatal(err)
			}

			if rec.Code != tc.wantCode {
				t.Errorf("ステータスコード不一致: got=%d, want=%d", rec.Code, tc.wantCode)
			}
			if !strings.Contains(rec.Body.String(), tc.wantBody) {
				t.Errorf("期待するメッセージが含まれていない: want=%q, got=%q", tc.wantBody, rec.Body.String())
			}
		})
	}
}
