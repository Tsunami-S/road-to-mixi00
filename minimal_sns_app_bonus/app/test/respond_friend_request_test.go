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

func setupTestDB_Respond(t *testing.T) {
	db.DB = initTestDB()
}

func TestRespondFriendRequest_Scenarios(t *testing.T) {
	setupTestDB_Respond(t)
	e := echo.New()

	handler := create.NewFriendRespondHandler(
		&validate.RealValidator{},
		&repo_create.RealFriendRespondRepository{},
	)

	tests := []struct {
		name     string
		body     string
		wantCode int
		wantBody string
	}{
		{"1.正常に申請を承認", `{"user1_id":"id45", "user2_id":"id1", "action":"accepted"}`, http.StatusOK, "request accepted"},
		{"2.正常に申請を拒否", `{"user1_id":"id11", "user2_id":"id27", "action":"rejected"}`, http.StatusOK, "request rejected"},
		{"3.自分自身に応答", `{"user1_id":"id1", "user2_id":"id1", "action":"accepted"}`, http.StatusBadRequest, "invalid user IDs"},
		{"4.存在しない申請", `{"user1_id":"id3", "user2_id":"id4", "action":"accepted"}`, http.StatusBadRequest, "request not found"},
		{"5.不正なアクション", `{"user1_id":"id45", "user2_id":"id1", "action":"maybe"}`, http.StatusBadRequest, "invalid action"},
		{"6.無効な user1_id", `{"user1_id":"invalid", "user2_id":"id1", "action":"accepted"}`, http.StatusBadRequest, "user not found"},
		{"7.無効な user2_id", `{"user1_id":"id1", "user2_id":"invalid", "action":"accepted"}`, http.StatusBadRequest, "user not found"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/respond_friend_request", strings.NewReader(tc.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if err := handler.RespondRequest(c); err != nil {
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
