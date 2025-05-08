package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"minimal_sns_app/db"
	"minimal_sns_app/handler/create"
	"minimal_sns_app/handler/validate"
	repo_create "minimal_sns_app/repository/create"

	"github.com/labstack/echo/v4"
)

func setupTestDB_Request(t *testing.T) {
	db.DB = initTestDB()
}

func TestRequestFriend_Scenarios(t *testing.T) {
	setupTestDB_Request(t)
	e := echo.New()

	handler := create.NewRequestFriendHandler(
		&validate.RealValidator{},
		&repo_create.RealFriendRequestRepository{},
	)

	tests := []struct {
		name     string
		body     string
		wantCode int
		wantBody string
	}{
		{"1.正常なフレンド申請", `{"user1_id":"id32","user2_id":"id43"}`, http.StatusOK, "friend request sent"},
		{"2.自分自身に申請", `{"user1_id":"id1","user2_id":"id1"}`, http.StatusBadRequest, "cannot request yourself"},
		{"3.ブロックしている", `{"user1_id":"id1","user2_id":"id39"}`, http.StatusBadRequest, "cannot send friend request due to block"},
		{"4.ブロックされている", `{"user1_id":"id1","user2_id":"id40"}`, http.StatusBadRequest, "cannot send friend request due to block"},
		{"5.すでにフレンド", `{"user1_id":"id1","user2_id":"id2"}`, http.StatusBadRequest, "you are already friends"},
		{"6.すでにフレンド(相手から)", `{"user1_id":"id1","user2_id":"id4"}`, http.StatusBadRequest, "you are already friends"},
		{"7.逆方向にpendingな申請がある", `{"user1_id":"id41","user2_id":"id1"}`, http.StatusBadRequest, "you already have a pending request from this user"},
		{"8.同じ方向の申請がすでにある", `{"user1_id":"id1","user2_id":"id27"}`, http.StatusBadRequest, "friend request already sent"},
		{"9.存在しない user1_id", `{"user1_id":"invalid_user","user2_id":"id2"}`, http.StatusBadRequest, "user not found"},
		{"10.存在しない user2_id", `{"user1_id":"id2","user2_id":"invalid_user"}`, http.StatusBadRequest, "user not found"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/request_friend", strings.NewReader(tc.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

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
