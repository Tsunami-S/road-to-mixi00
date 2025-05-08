package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"minimal_sns_app/db"
	"minimal_sns_app/handler/get"
	"minimal_sns_app/handler/validate"
	"minimal_sns_app/interfaces"
	repo_get "minimal_sns_app/repository/get"

	"github.com/labstack/echo/v4"
)

func setupTestDB_FOF_Paging(t *testing.T) {
	db.DB = initTestDB()
}

func TestGetFriendOfFriendListPaging_Scenarios(t *testing.T) {
	setupTestDB_FOF_Paging(t)
	e := echo.New()

	handler := get.NewFriendOfFriendPagingHandler(
		&validate.RealValidator{},
		&validate.RealPaginationValidator{},
		&repo_get.RealFriendOfFriendPagingRepository{},
	)

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
			name:     "1.ページ1に特定のユーザーが含まれる",
			userID:   "id1",
			limit:    "2",
			page:     "1",
			wantCode: http.StatusOK,
			wantBody: "user13",
		},
		{
			name:     "2.ページ2に他のユーザーが出現する",
			userID:   "id1",
			limit:    "2",
			page:     "2",
			wantCode: http.StatusOK,
			wantBody: "user12",
		},
		{
			name:     "3.最終ページはデータがない",
			userID:   "id1",
			limit:    "2",
			page:     "99",
			wantCode: http.StatusOK,
			wantBody: "no friends of friends found",
		},
		{
			name:     "4.存在しないID",
			userID:   "invalid_id",
			limit:    "2",
			page:     "1",
			wantCode: http.StatusBadRequest,
			wantBody: "user ID not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url := "/get_friend_of_friend_list_paging?id=" + tc.userID + "&limit=" + tc.limit + "&page=" + tc.page
			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if err := handler.FriendOfFriendPaging(c); err != nil {
				t.Fatal(err)
			}

			body := rec.Body.String()

			if rec.Code != tc.wantCode {
				t.Errorf("ステータスコード不一致: got %d, want %d", rec.Code, tc.wantCode)
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
