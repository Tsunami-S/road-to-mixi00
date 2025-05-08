package validate

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"minimal_sns_app/db"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) {
	db.InitDB()
}

func TestUserExists(t *testing.T) {
	setupTestDB(t)

	tests := []struct {
		name     string
		id       int
		expected bool
	}{
		{
			name:     "存在するユーザーID",
			id:       1,
			expected: true,
		},
		{
			name:     "存在しないユーザーID",
			id:       9999,
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := UserExists(tc.id)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestParseAndValidatePagination(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name      string
		limit     string
		page      string
		wantLimit int
		wantPage  int
		wantErr   string
	}{
		{
			name:      "正常なページとリミット",
			limit:     "10",
			page:      "2",
			wantLimit: 10,
			wantPage:  2,
			wantErr:   "",
		},
		{
			name:    "limit が非数値",
			limit:   "abc",
			page:    "1",
			wantErr: "error: invalid limit",
		},
		{
			name:    "limit が 0",
			limit:   "0",
			page:    "1",
			wantErr: "error: invalid limit",
		},
		{
			name:    "page が非数値",
			limit:   "5",
			page:    "xyz",
			wantErr: "error: invalid page",
		},
		{
			name:    "page が負数",
			limit:   "5",
			page:    "-1",
			wantErr: "error: invalid page",
		},
		{
			name:    "両方不正（limit 優先で検出）",
			limit:   "aaa",
			page:    "bbb",
			wantErr: "error: invalid limit",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/?limit="+tc.limit+"&page="+tc.page, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			limit, page, err := (&RealPaginationValidator{}).ParseAndValidatePagination(c)

			if tc.wantErr != "" {
				assert.Error(t, err)
				httpErr, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tc.wantErr, httpErr.Message)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantLimit, limit)
				assert.Equal(t, tc.wantPage, page)
			}
		})
	}
}
