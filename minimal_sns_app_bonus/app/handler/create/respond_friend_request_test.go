package create_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"minimal_sns_app/handler/create"
	"minimal_sns_app/model"
	"minimal_sns_app/test/mock"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRespondRequest(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name         string
		input        model.RespondRequestInput
		validator    *mock.Validator
		repo         *mock.MockRespondRepo
		wantCode     int
		wantContains string
	}{
		{
			name: "✅ success: accepted",
			input: model.RespondRequestInput{
				User1ID: "user01", User2ID: "user02", Action: "accepted",
			},
			validator:    &mock.Validator{Exists: true},
			repo:         &mock.MockRespondRepo{},
			wantCode:     http.StatusOK,
			wantContains: "request accepted",
		},
		{
			name: "✅ success: rejected",
			input: model.RespondRequestInput{
				User1ID: "user01", User2ID: "user02", Action: "rejected",
			},
			validator:    &mock.Validator{Exists: true},
			repo:         &mock.MockRespondRepo{},
			wantCode:     http.StatusOK,
			wantContains: "request rejected",
		},
		{
			name: "❌ invalid user ID (empty)",
			input: model.RespondRequestInput{
				User1ID: "", User2ID: "user02", Action: "accepted",
			},
			validator:    &mock.Validator{Exists: true},
			repo:         &mock.MockRespondRepo{},
			wantCode:     http.StatusBadRequest,
			wantContains: "invalid user IDs",
		},
		{
			name: "❌ invalid action",
			input: model.RespondRequestInput{
				User1ID: "user01", User2ID: "user02", Action: "maybe",
			},
			validator:    &mock.Validator{Exists: true},
			repo:         &mock.MockRespondRepo{},
			wantCode:     http.StatusBadRequest,
			wantContains: "invalid action",
		},
		{
			name: "❌ same user",
			input: model.RespondRequestInput{
				User1ID: "user01", User2ID: "user01", Action: "accepted",
			},
			validator:    &mock.Validator{Exists: true},
			repo:         &mock.MockRespondRepo{},
			wantCode:     http.StatusBadRequest,
			wantContains: "invalid user IDs",
		},
		{
			name: "❌ user1 not found",
			input: model.RespondRequestInput{
				User1ID: "notfound", User2ID: "user02", Action: "accepted",
			},
			validator:    &mock.Validator{Exists: false},
			repo:         &mock.MockRespondRepo{},
			wantCode:     http.StatusBadRequest,
			wantContains: "user ID not found",
		},
		{
			name: "❌ user2 not found",
			input: model.RespondRequestInput{
				User1ID: "user01", User2ID: "none", Action: "accepted",
			},
			validator:    &mock.Validator{Exists: false},
			repo:         &mock.MockRespondRepo{},
			wantCode:     http.StatusBadRequest,
			wantContains: "user ID not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.input)
			req := httptest.NewRequest(http.MethodPost, "/respond_friend_request", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := create.NewFriendRespondHandler(tc.validator, tc.repo)
			err := handler.RespondRequest(c)

			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantContains)
		})
	}
}
