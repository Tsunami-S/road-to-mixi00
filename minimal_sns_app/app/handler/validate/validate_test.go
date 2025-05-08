package validate

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestParseAndValidatePagination(t *testing.T) {
	tests := []struct {
		name      string
		limit     string
		page      string
		wantLimit int
		wantPage  int
		wantErr   string
	}{
		{
			name:      "正常系: limit=10, page=2",
			limit:     "10",
			page:      "2",
			wantLimit: 10,
			wantPage:  2,
			wantErr:   "",
		},
		{
			name:    "異常系: limitが非数値",
			limit:   "abc",
			page:    "2",
			wantErr: "error: invalid limit",
		},
		{
			name:    "異常系: limitが0",
			limit:   "0",
			page:    "2",
			wantErr: "error: invalid limit",
		},
		{
			name:    "異常系: pageが非数値",
			limit:   "10",
			page:    "abc",
			wantErr: "error: invalid page",
		},
		{
			name:    "異常系: pageが負数",
			limit:   "10",
			page:    "-1",
			wantErr: "error: invalid page",
		},
		{
			name:    "異常系: 両方不正（limit優先）",
			limit:   "xyz",
			page:    "0",
			wantErr: "error: invalid limit",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/?limit="+tc.limit+"&page="+tc.page, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			limit, page, err := ParseAndValidatePagination(c)

			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("エラーが期待されたが nil")
				}
				httpErr, ok := err.(*echo.HTTPError)
				if !ok || httpErr.Message != tc.wantErr {
					t.Errorf("期待されたエラー %q, 実際のエラー: %v", tc.wantErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("期待されないエラー: %v", err)
				}
				if limit != tc.wantLimit || page != tc.wantPage {
					t.Errorf("limit, page 値不一致: got (%d, %d), want (%d, %d)",
						limit, page, tc.wantLimit, tc.wantPage)
				}
			}
		})
	}
}
