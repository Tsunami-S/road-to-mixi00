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

func setupTestDB_Request(t *testing.T) {
	db.DB = initTestDB()
}

func TestRequestFriend_Scenarios(t *testing.T) {
	setupTestDB_Request(t)
	e := echo.New()

	tests := []struct {
		name     string
		user1ID  string
		user2ID  string
		wantCode int
		wantBody string
	}{
		{
			name:     "1.正常なフレンド申請",
			user1ID:  "id32",
			user2ID:  "id43",
			wantCode: http.StatusOK,
			wantBody: "friend request sent",
		},
		{
			name:     "2.自分自身に申請",
			user1ID:  "id1",
			user2ID:  "id1",
			wantCode: http.StatusBadRequest,
			wantBody: "cannot request yourself",
		},
		{
			name:     "3.ブロックしている",
			user1ID:  "id1",
			user2ID:  "id39",
			wantCode: http.StatusBadRequest,
			wantBody: "cannot send friend request due to block",
		},
		{
			name:     "4.ブロックされている",
			user1ID:  "id1",
			user2ID:  "id40",
			wantCode: http.StatusBadRequest,
			wantBody: "cannot send friend request due to block",
		},
		{
			name:     "5.すでにフレンド",
			user1ID:  "id1",
			user2ID:  "id2",
			wantCode: http.StatusBadRequest,
			wantBody: "you are already friends",
		},
		{
			name:     "6.すでにフレンド(相手から)",
			user1ID:  "id1",
			user2ID:  "id4",
			wantCode: http.StatusBadRequest,
			wantBody: "you are already friends",
		},
		{
			name:     "7.逆方向にpendingな申請がある",
			user1ID:  "id41",
			user2ID:  "id1", 
			wantCode: http.StatusBadRequest,
			wantBody: "you already have a pending request from this user",
		},
		{
			name:     "8.同じ方向の申請がすでにある",
			user1ID:  "id1",
			user2ID:  "id27",
			wantCode: http.StatusBadRequest,
			wantBody: "friend request already sent",
		},
		{
			name:     "9.存在しない user1_id",
			user1ID:  "invalid_user",
			user2ID:  "id2",
			wantCode: http.StatusBadRequest,
			wantBody: "user1_id: user ID not found",
		},
		{
			name:     "10.存在しない user2_id",
			user1ID:  "id2",
			user2ID:  "invalid_user",
			wantCode: http.StatusBadRequest,
			wantBody: "user2_id: user ID not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url := "/request_friend?user1_id=" + tc.user1ID + "&user2_id=" + tc.user2ID
			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if err := handler.RequestFriend(c); err != nil {
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
