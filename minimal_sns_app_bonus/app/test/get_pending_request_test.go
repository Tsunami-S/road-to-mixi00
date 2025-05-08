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

func setupTestDB_Pending(t *testing.T) {
	db.DB = initTestDB()
}

func TestGetPendingRequests_Scenarios(t *testing.T) {
	setupTestDB_Pending(t)
	e := echo.New()

	handler := get.NewPendingRequestHandler(
		&validate.RealValidator{},
		&repo_get.RealFriendRequestRepository{},
	)

	tests := []struct {
		name      string
		userID    string
		wantCode  int
		wantBody  string
		notInBody string
	}{
		{
			name:     "1.id41 に保留中リクエストあり",
			userID:   "id41",
			wantCode: http.StatusOK,
			wantBody: `"user1_id":"id1"`,
		},
		{
			name:     "2.id46 にリクエスト（rejected なので含まれない）",
			userID:   "id46",
			wantCode: http.StatusOK,
			wantBody: "no pending requests found",
		},
		{
			name:     "3.リクエストが 0 件（id5）",
			userID:   "id5",
			wantCode: http.StatusOK,
			wantBody: "no pending requests found",
		},
		{
			name:     "4.無効な ID",
			userID:   "invalid_id",
			wantCode: http.StatusBadRequest,
			wantBody: "user not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/pending_requests?user_id="+tc.userID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if err := handler.PendingRequests(c); err != nil {
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
		})
	}
}
