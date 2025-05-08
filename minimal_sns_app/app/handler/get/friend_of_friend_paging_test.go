package get

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"minimal_sns_app/interfaces"
	"minimal_sns_app/model"
	"minimal_sns_app/test/mock"

	"github.com/labstack/echo/v4"
)

func TestFriendOfFriendPagingHandler(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name      string
		id        string
		validator interfaces.UserValidator
		pager     interfaces.PaginationValidator
		repo      interfaces.FriendOfFriendPagingRepository
		wantCode  int
		wantBody  string
	}{
		{
			name:      "正常系: ページング取得成功",
			id:        "1",
			validator: &mock.UserValidatorMock{Exist: true},
			pager:     &mock.PaginationValidatorMock{Limit: 2, Page: 1},
			repo: &mock.FriendOfFriendPagingRepoMock{
				Result: []model.Friend{{ID: 2, Name: "user02"}},
			},
			wantCode: http.StatusOK,
			wantBody: "user02",
		},
		{
			name:      "異常系: IDが空文字",
			id:        "",
			validator: &mock.UserValidatorMock{},
			pager:     &mock.PaginationValidatorMock{},
			repo:      &mock.FriendOfFriendPagingRepoMock{},
			wantCode:  http.StatusBadRequest,
			wantBody:  "id is required",
		},
		{
			name:      "異常系: IDが数値でない",
			id:        "abc",
			validator: &mock.UserValidatorMock{},
			pager:     &mock.PaginationValidatorMock{},
			repo:      &mock.FriendOfFriendPagingRepoMock{},
			wantCode:  http.StatusBadRequest,
			wantBody:  "id must be a positive integer",
		},
		{
			name:      "異常系: IDが0以下",
			id:        "0",
			validator: &mock.UserValidatorMock{},
			pager:     &mock.PaginationValidatorMock{},
			repo:      &mock.FriendOfFriendPagingRepoMock{},
			wantCode:  http.StatusBadRequest,
			wantBody:  "id must be a positive integer",
		},
		{
			name:      "異常系: ユーザー存在チェック中にDBエラー",
			id:        "1",
			validator: &mock.UserValidatorMock{Err: echo.NewHTTPError(http.StatusInternalServerError, "DB error")},
			pager:     &mock.PaginationValidatorMock{Limit: 2, Page: 1},
			repo:      &mock.FriendOfFriendPagingRepoMock{},
			wantCode:  http.StatusInternalServerError,
			wantBody:  "DB error",
		},
		{
			name:      "異常系: ユーザーが存在しない",
			id:        "9999",
			validator: &mock.UserValidatorMock{Exist: false},
			pager:     &mock.PaginationValidatorMock{Limit: 2, Page: 1},
			repo:      &mock.FriendOfFriendPagingRepoMock{},
			wantCode:  http.StatusBadRequest,
			wantBody:  "user not found",
		},
		{
			name:      "異常系: リポジトリでDBエラー発生",
			id:        "1",
			validator: &mock.UserValidatorMock{Exist: true},
			pager:     &mock.PaginationValidatorMock{Limit: 2, Page: 1},
			repo:      &mock.FriendOfFriendPagingRepoMock{Err: echo.NewHTTPError(http.StatusInternalServerError, "DB query error")},
			wantCode:  http.StatusInternalServerError,
			wantBody:  "DB query error",
		},
		{
			name:      "正常系: 友達の友達が存在しない",
			id:        "1",
			validator: &mock.UserValidatorMock{Exist: true},
			pager:     &mock.PaginationValidatorMock{Limit: 10, Page: 1},
			repo:      &mock.FriendOfFriendPagingRepoMock{Result: []model.Friend{}},
			wantCode:  http.StatusOK,
			wantBody:  "no friends of friends found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/get_friend_of_friend_list_paging?id="+tc.id, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := NewFriendOfFriendPagingHandler(tc.validator, tc.pager, tc.repo)
			if err := handler.FriendOfFriendPaging(c); err != nil {
				t.Fatal(err)
			}

			if rec.Code != tc.wantCode {
				t.Errorf("ステータスコード不一致: got %d, want %d", rec.Code, tc.wantCode)
			}
			if !strings.Contains(rec.Body.String(), tc.wantBody) {
				t.Errorf("期待する文字列が含まれない: want %q, got %q", tc.wantBody, rec.Body.String())
			}
		})
	}
}
