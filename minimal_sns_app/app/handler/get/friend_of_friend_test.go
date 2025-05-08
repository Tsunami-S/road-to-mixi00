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

func TestFriendOfFriendHandler(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name      string
		id        string
		validator interfaces.UserValidator
		repo      interfaces.FriendOfFriendRepository
		wantCode  int
		wantBody  string
	}{
		{
			name:      "正常系: 友達の友達が存在",
			id:        "1",
			validator: &mock.UserValidatorMock{Exist: true},
			repo: &mock.FriendOfFriendRepoMock{
				Result: []model.Friend{{ID: 2, Name: "user02"}},
			},
			wantCode: http.StatusOK,
			wantBody: "user02",
		},
		{
			name:      "異常系: idが空",
			id:        "",
			validator: &mock.UserValidatorMock{},
			repo:      &mock.FriendOfFriendRepoMock{},
			wantCode:  http.StatusBadRequest,
			wantBody:  "id is required",
		},
		{
			name:      "異常系: 数値以外",
			id:        "abc",
			validator: &mock.UserValidatorMock{},
			repo:      &mock.FriendOfFriendRepoMock{},
			wantCode:  http.StatusBadRequest,
			wantBody:  "id must be a positive integer",
		},
		{
			name:      "異常系: 存在しないユーザー",
			id:        "9999",
			validator: &mock.UserValidatorMock{Exist: false},
			repo:      &mock.FriendOfFriendRepoMock{},
			wantCode:  http.StatusBadRequest,
			wantBody:  "user not found",
		},
		{
			name:      "正常系: 友達の友達がいない",
			id:        "1",
			validator: &mock.UserValidatorMock{Exist: true},
			repo:      &mock.FriendOfFriendRepoMock{Result: []model.Friend{}},
			wantCode:  http.StatusOK,
			wantBody:  "no friends of friends found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/get_friend_of_friend_list?id="+tc.id, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := NewFriendOfFriendHandler(tc.validator, tc.repo)
			if err := handler.FriendOfFriend(c); err != nil {
				t.Fatal(err)
			}

			body := rec.Body.String()
			if rec.Code != tc.wantCode {
				t.Errorf("ステータスコード不一致: got %d, want %d", rec.Code, tc.wantCode)
			}
			if !strings.Contains(body, tc.wantBody) {
				t.Errorf("期待する文字列が含まれない: want=%q, got=%q", tc.wantBody, body)
			}
		})
	}
}
