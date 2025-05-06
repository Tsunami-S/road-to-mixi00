package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetFriendOfFriendListPaging_Scenarios(t *testing.T) {
	setupTestDB_FOF(t)
	e := echo.New()

	tests := []struct {
		name      string
		userID    string
		limit     string
		page      string
		wantCode  int
		wantBody  string
		notInBody string
	}{
		{
			name:     "✅ ページ1に特定のユーザーが含まれる",
			userID:   "1",
			limit:    "2",
			page:     "1",
			wantCode: http.StatusOK,
			wantBody: "user11",
		},
		{
			name:     "✅ ページ2に他のユーザーが出現する",
			userID:   "1",
			limit:    "2",
			page:     "2",
			wantCode: http.StatusOK,
			wantBody: "user13",
		},
		{
			name:     "🟩 最終ページはデータがない",
			userID:   "1",
			limit:    "2",
			page:     "99",
			wantCode: http.StatusOK,
			wantBody: "no friends of friends found",
		},
		{
			name:     "❌ 存在しないID",
			userID:   "9999",
			limit:    "2",
			page:     "1",
			wantCode: http.StatusBadRequest,
			wantBody: "user not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url := "/get_friend_of_friend_list_paging?id=" + tc.userID + "&limit=" + tc.limit + "&page=" + tc.page
			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if err := GetFriendOfFriendListPaging(c); err != nil {
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
		})
	}
}
