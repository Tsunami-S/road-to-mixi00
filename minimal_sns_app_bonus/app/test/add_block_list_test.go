package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"minimal_sns_app/db"
	"minimal_sns_app/handler/create"

	"github.com/labstack/echo/v4"
)

func setupTestDB_Block(t *testing.T) {
	db.DB = initTestDB()
}

func TestAddBlockList_Scenarios(t *testing.T) {
	setupTestDB_Block(t)
	e := echo.New()

	tests := []struct {
		name     string
		body     string
		wantCode int
		wantBody string
	}{
		{
			name:     "1.ブロック成功（フレンド削除・申請拒否）",
			body:     `{"user1_id":"id1", "user2_id":"id44"}`,
			wantCode: http.StatusOK,
			wantBody: "user blocked",
		},
		{
			name:     "2.自分自身をブロック",
			body:     `{"user1_id":"id1", "user2_id":"id1"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "cannot block yourself",
		},
		{
			name:     "3.既にブロック済み（user1 -> user2）",
			body:     `{"user1_id":"id1", "user2_id":"id39"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "already blocked",
		},
		{
			name:     "4.既にブロック済み（user2 -> user1）",
			body:     `{"user1_id":"id38", "user2_id":"id1"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "already blocked",
		},
		{
			name:     "5.存在しない user1_id",
			body:     `{"user1_id":"invalid_id", "user2_id":"id2"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "user1_id: user ID not found",
		},
		{
			name:     "6.存在しない user2_id",
			body:     `{"user1_id":"id2", "user2_id":"invalid_id"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "user2_id: user ID not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/block_user", strings.NewReader(tc.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if err := create.AddBlockList(c); err != nil {
				t.Fatal(err)
			}

			body := rec.Body.String()

			if rec.Code != tc.wantCode {
				t.Errorf("ステータスコード不一致: got=%d, want=%d", rec.Code, tc.wantCode)
			}
			if !strings.Contains(body, tc.wantBody) {
				t.Errorf("期待する文字列が含まれていない: want=%q, got=%q", tc.wantBody, body)
			}
		})
	}
}
