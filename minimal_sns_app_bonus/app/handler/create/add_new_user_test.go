package create

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"minimal_sns_app/model"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockUserRepo struct {
	exists    bool
	existsErr error
	createErr error
}

func (m *mockUserRepo) IsUserIDExists(id string) (bool, error) {
	return m.exists, m.existsErr
}

func (m *mockUserRepo) CreateUser(id, name string) error {
	return m.createErr
}

func TestAddNewUser(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name         string
		input        model.AddUserRequest
		repo         *mockUserRepo
		wantCode     int
		wantContains string
	}{
		{
			name: "success: user added",
			input: model.AddUserRequest{
				ID:   "user01",
				Name: "Taro",
			},
			repo:         &mockUserRepo{},
			wantCode:     http.StatusOK,
			wantContains: "user added",
		},
		{
			name: "error: empty ID",
			input: model.AddUserRequest{
				ID:   "",
				Name: "Taro",
			},
			repo:         &mockUserRepo{},
			wantCode:     http.StatusBadRequest,
			wantContains: "id must have",
		},
		{
			name: "error: empty Name",
			input: model.AddUserRequest{
				ID:   "user01",
				Name: "",
			},
			repo:         &mockUserRepo{},
			wantCode:     http.StatusBadRequest,
			wantContains: "name must have",
		},
		{
			name: "error: ID too long",
			input: model.AddUserRequest{
				ID:   "abcdefghijklmnopqrstuvwxyz",
				Name: "Taro",
			},
			repo:         &mockUserRepo{},
			wantCode:     http.StatusBadRequest,
			wantContains: "id must have",
		},
		{
			name: "error: Name too long",
			input: model.AddUserRequest{
				ID:   "user01",
				Name: string(make([]byte, 65)), // 65文字
			},
			repo:         &mockUserRepo{},
			wantCode:     http.StatusBadRequest,
			wantContains: "name must have",
		},
		{
			name: "error: duplicate ID",
			input: model.AddUserRequest{
				ID:   "user01",
				Name: "Taro",
			},
			repo:         &mockUserRepo{exists: true},
			wantCode:     http.StatusBadRequest,
			wantContains: "already exists",
		},
		{
			name: "error: exists check fails",
			input: model.AddUserRequest{
				ID:   "user01",
				Name: "Taro",
			},
			repo:         &mockUserRepo{existsErr: errors.New("db error")},
			wantCode:     http.StatusInternalServerError,
			wantContains: "failed to check",
		},
		{
			name: "error: create fails",
			input: model.AddUserRequest{
				ID:   "user01",
				Name: "Taro",
			},
			repo:         &mockUserRepo{createErr: errors.New("insert failed")},
			wantCode:     http.StatusInternalServerError,
			wantContains: "failed to create",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			payload, _ := json.Marshal(tc.input)
			req := httptest.NewRequest(http.MethodPost, "/add_new_user", bytes.NewReader(payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := NewUserHandler(tc.repo)
			err := h.AddNewUser(c)

			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantContains)
		})
	}
}
