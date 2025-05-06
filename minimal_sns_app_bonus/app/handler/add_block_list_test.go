package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/db"
	"minimal_sns_app/test"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestDB_Block(t *testing.T) {
	db.DB = test.InitTestDB()
}

func TestAddBlockList_Scenarios(t *testing.T) {
	setupTestDB_Block(t)
	e := echo.New()

	tests := []struct {
		name     string
		user1ID  string
		user2ID  string
		wantCode int
		wantBody string
	}{
		{
			name:     "✅ ブロック成功（フレンド削除・申請拒否）",
			user1ID:  "id1",
			user2ID:  "id44", // pending中のフレンド申請あり
			wantCode: http.StatusOK,
			wantBody: "user blocked",
		},
		{
			name:     "❌ 自分自身をブロック",
			user1ID:  "id1",
			user2ID:  "id1",
			wantCode: http.StatusBadRequest,
			wantBody: "cannot block yourself",
		},
		{
			name:     "❌ 既にブロック済み",
			user1ID:  "id1",
			user2ID:  "id39", // test.sql で事前にブロック関係
			wantCode: http.StatusBadRequest,
			wantBody: "already blocked",
		},
		{
			name:     "❌ 存在しない user1_id",
			user1ID:  "invalid_id",
			user2ID:  "id2",
			wantCode: http.StatusBadRequest,
			wantBody: "user1_id: user ID not found",
		},
		{
			name:     "❌ 存在しない user2_id",
			user1ID:  "id2",
			user2ID:  "invalid_id",
			wantCode: http.StatusBadRequest,
			wantBody: "user2_id: user ID not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url := "/block_user?user1_id=" + tc.user1ID + "&user2_id=" + tc.user2ID
			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if err := AddBlockList(c); err != nil {
				t.Fatal(err)
			}

			if rec.Code != tc.wantCode {
				t.Errorf("ステータスコード不一致: got=%d, want=%d", rec.Code, tc.wantCode)
			}

			body := rec.Body.String()
			if !strings.Contains(body, tc.wantBody) {
				t.Errorf("期待する文字列が含まれていない: want=%q, got=%q", tc.wantBody, body)
			}
		})
	}
}
