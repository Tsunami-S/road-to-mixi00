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

func TestFriendHandler(t *testing.T) {
	e := echo.New()
	tests := []struct {
		name      string
		id        string
		validator interfaces.UserValidator
		repo      interfaces.FriendRepository
		wantCode  int
		wantBody  string
	}{
		{
			name:      "正常系: フレンドあり",
			id:        "1",
			validator: &mock.UserValidatorMock{Exist: true},
			repo:      &mock.FriendRepoMock{Result: []model.Friend{{ID: 2, Name: "user02"}}},
			wantCode:  http.StatusOK,
			wantBody:  "user02",
		},
		{
			name:      "異常系: IDが空",
			id:        "",
			validator: &mock.UserValidatorMock{},
			repo:      &mock.FriendRepoMock{},
			wantCode:  http.StatusBadRequest,
			wantBody:  "id is required",
		},
		{
			name:      "異常系: 存在しない",
			id:        "9999",
			validator: &mock.UserValidatorMock{Exist: false},
			repo:      &mock.FriendRepoMock{},
			wantCode:  http.StatusBadRequest,
			wantBody:  "user not found",
		},
		{
			name:      "正常系: フレンドなし",
			id:        "1",
			validator: &mock.UserValidatorMock{Exist: true},
			repo:      &mock.FriendRepoMock{Result: []model.Friend{}},
			wantCode:  http.StatusOK,
			wantBody:  "no friends found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/get_friend_list?id="+tc.id, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := NewFriendHandler(tc.validator, tc.repo)
			if err := h.Friend(c); err != nil {
				t.Fatal(err)
			}

			if rec.Code != tc.wantCode {
				t.Errorf("want %d, got %d", tc.wantCode, rec.Code)
			}
			if !strings.Contains(rec.Body.String(), tc.wantBody) {
				t.Errorf("body missing %q, got %q", tc.wantBody, rec.Body.String())
			}
		})
	}
}
