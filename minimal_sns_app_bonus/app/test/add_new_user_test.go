package test

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/db"
	"minimal_sns_app/handler/create"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestDB_AddUser(t *testing.T) {
	db.DB = initTestDB()
}

func TestAddNewUser_Scenarios(t *testing.T) {
	setupTestDB_AddUser(t)
	e := echo.New()

	tests := []struct {
		name     string
		id       string
		userName string
		wantCode int
		wantBody string
	}{
		{
			name:     "1.ユーザー追加成功",
			id:       "new_user_1",
			userName: "テストユーザー",
			wantCode: http.StatusOK,
			wantBody: "user added",
		},
		{
			name:     "2.同じIDのユーザーが既に存在する",
			id:       "id1",
			userName: "新しい名前",
			wantCode: http.StatusBadRequest,
			wantBody: "user ID already exists",
		},
		{
			name:     "3.ID が空",
			id:       "",
			userName: "valid",
			wantCode: http.StatusBadRequest,
			wantBody: "id must have 1 ~ 20 characters",
		},
		{
			name:     "4.ID が長すぎる",
			id:       strings.Repeat("a", 21),
			userName: "valid",
			wantCode: http.StatusBadRequest,
			wantBody: "id must have 1 ~ 20 characters",
		},
		{
			name:     "5.名前が空",
			id:       "valid_id",
			userName: "",
			wantCode: http.StatusBadRequest,
			wantBody: "name must have 1 ~ 64 characters",
		},
		{
			name:     "6.名前が長すぎる",
			id:       "valid_id",
			userName: strings.Repeat("あ", 65),
			wantCode: http.StatusBadRequest,
			wantBody: "name must have 1 ~ 64 characters",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url := "/add_new_user?id=" + tc.id + "&name=" + tc.userName
			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := create.AddNewUser(c)
			if err != nil {
				t.Fatal(err)
			}

			if rec.Code != tc.wantCode {
				t.Errorf("ステータスコード不一致: got=%d, want=%d", rec.Code, tc.wantCode)
			}
			if !strings.Contains(rec.Body.String(), tc.wantBody) {
				t.Errorf("期待する文字列が含まれない: want=%q, got=%q", tc.wantBody, rec.Body.String())
			}
		})
	}
}
