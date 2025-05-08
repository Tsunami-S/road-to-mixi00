package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"minimal_sns_app/db"
	"minimal_sns_app/handler/create"
	repo_create "minimal_sns_app/repository/create"

	"github.com/labstack/echo/v4"
)

func setupTestDB_AddUser(t *testing.T) {
	db.DB = initTestDB()
}

func TestAddNewUser_Scenarios(t *testing.T) {
	setupTestDB_AddUser(t)
	e := echo.New()

	handler := create.NewAddUserHandler(&repo_create.RealUserRepository{})

	tests := []struct {
		name     string
		body     string
		wantCode int
		wantBody string
	}{
		{
			name:     "1.ユーザー追加成功",
			body:     `{"id":"new_user_1", "name":"テストユーザー"}`,
			wantCode: http.StatusOK,
			wantBody: "user added",
		},
		{
			name:     "2.同じIDのユーザーが既に存在する",
			body:     `{"id":"id1", "name":"新しい名前"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "user ID already exists",
		},
		{
			name:     "3.ID が空",
			body:     `{"id":"", "name":"valid"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "id must have 1 ~ 20 characters",
		},
		{
			name:     "4.ID が長すぎる",
			body:     `{"id":"` + strings.Repeat("a", 21) + `", "name":"valid"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "id must have 1 ~ 20 characters",
		},
		{
			name:     "5.名前が空",
			body:     `{"id":"valid_id", "name":""}`,
			wantCode: http.StatusBadRequest,
			wantBody: "name must have 1 ~ 64 characters",
		},
		{
			name:     "6.名前が長すぎる",
			body:     `{"id":"valid_id", "name":"` + strings.Repeat("あ", 65) + `"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "name must have 1 ~ 64 characters",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/add_new_user", strings.NewReader(tc.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.AddNewUser(c)
			if err != nil {
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
