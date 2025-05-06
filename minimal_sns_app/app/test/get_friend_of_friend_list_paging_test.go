package test

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/handler"
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
		{"1.ページ1に特定のユーザーが含まれる", "1", "2", "1", http.StatusOK, "user11", ""},
		{"2.ページ2に他のユーザーが出現する", "1", "2", "2", http.StatusOK, "user13", ""},
		{"3.最終ページはデータがない", "1", "2", "99", http.StatusOK, "no friends of friends found", ""},
		{"4.存在しないID", "9999", "2", "1", http.StatusBadRequest, "user not found", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url := "/get_friend_of_friend_list_paging?id=" + tc.userID + "&limit=" + tc.limit + "&page=" + tc.page
			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if err := handler.GetFriendOfFriendListPaging(c); err != nil {
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

