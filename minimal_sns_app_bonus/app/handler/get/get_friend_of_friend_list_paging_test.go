package get

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"minimal_sns_app/model"
	"minimal_sns_app/test/mock"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestFriendOfFriendPagingHandler(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name           string
		userID         string
		limit, page    int
		valExists      bool
		valErr         error
		pageErr        error
		repoResult     []model.Friend
		repoErr        error
		wantCode       int
		wantBodyString string
	}{
		{
			name:           "正常系: ページ1でデータあり",
			userID:         "user01",
			limit:          2,
			page:           1,
			valExists:      true,
			repoResult:     []model.Friend{{ID: "u2", Name: "user02"}},
			wantCode:       http.StatusOK,
			wantBodyString: "user02",
		},
		{
			name:           "正常系: データなし",
			userID:         "user01",
			limit:          2,
			page:           99,
			valExists:      true,
			repoResult:     []model.Friend{},
			wantCode:       http.StatusOK,
			wantBodyString: "no friends of friends found",
		},
		{
			name:           "異常系: 存在しないユーザー",
			userID:         "unknown",
			valExists:      false,
			limit:          2,
			page:           1,
			wantCode:       http.StatusBadRequest,
			wantBodyString: "user not found",
		},
		{
			name:           "異常系: validator エラー",
			userID:         "userX",
			valErr:         errors.New("validator error"),
			limit:          2,
			page:           1,
			wantCode:       http.StatusBadRequest,
			wantBodyString: "validator error",
		},
		{
			name:           "異常系: ページングバリデーションエラー",
			userID:         "user01",
			valExists:      true,
			pageErr:        echo.NewHTTPError(http.StatusBadRequest, "invalid pagination"),
			wantCode:       http.StatusBadRequest,
			wantBodyString: "invalid pagination",
		},
		{
			name:           "異常系: repoエラー",
			userID:         "user01",
			valExists:      true,
			limit:          2,
			page:           1,
			repoErr:        errors.New("repo error"),
			wantCode:       http.StatusInternalServerError,
			wantBodyString: "repo error",
		},
		{
			name:           "異常系: 空ID",
			userID:         "",
			wantCode:       http.StatusBadRequest,
			wantBodyString: "user_id must be",
		},
		{
			name:           "異常系: ID長すぎ",
			userID:         strings.Repeat("a", 21),
			wantCode:       http.StatusBadRequest,
			wantBodyString: "user_id must be",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := &mock.UserValidatorMock{UserExistsResult: tc.valExists, Err: tc.valErr}
			p := &mock.PaginationValidatorMock{Limit: tc.limit, Page: tc.page, Err: tc.pageErr}
			r := &mock.FriendOfFriendPagingRepositoryMock{Result: tc.repoResult, Err: tc.repoErr}

			h := NewFriendOfFriendPagingHandler(v, p, r)

			req := httptest.NewRequest(http.MethodGet, "/?id="+tc.userID+"&limit=2&page=1", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.FriendOfFriendPaging(c)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.wantBodyString)
		})
	}
}
