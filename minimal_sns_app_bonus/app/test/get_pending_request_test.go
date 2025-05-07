package test

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/db"
	"minimal_sns_app/handler/get"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestDB_Pending(t *testing.T) {
	db.DB = initTestDB()
}

func TestGetPendingRequests_Scenarios(t *testing.T) {
	setupTestDB_Pending(t)
	e := echo.New()

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
			name:     "2.id46 にリクエスト（rejectedなので含まれない）",
			userID:   "id46",
			wantCode: http.StatusOK,
			wantBody: "no pending requests found",
		},
		{
			name:     "3.リクエストが0件（id5）",
			userID:   "id5",
			wantCode: http.StatusOK,
			wantBody: "no pending requests found",
		},
		{
			name:     "4.無効なID",
			userID:   "invalid_id",
			wantCode: http.StatusBadRequest,
			wantBody: "user ID not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url := "/pending_requests?user_id=" + tc.userID
			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := get.PendingRequests(c)
			if err != nil {
				t.Fatal(err)
			}

			if rec.Code != tc.wantCode {
				t.Errorf("ステータスコード不一致: got=%d, want=%d", rec.Code, tc.wantCode)
			}

			body := rec.Body.String()
			if tc.wantBody != "" && !strings.Contains(body, tc.wantBody) {
				t.Errorf("期待する文字列が含まれていない: want=%q, got=%q", tc.wantBody, body)
			}
			if tc.notInBody != "" && strings.Contains(body, tc.notInBody) {
				t.Errorf("含まれてはいけない文字列が含まれている: notWant=%q, got=%q", tc.notInBody, body)
			}
		})
	}
}
