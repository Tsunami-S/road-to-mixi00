package get

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"minimal_sns_app/interfaces"
)

type FriendOfFriendPagingHandler struct {
	Validator           interfaces.Validator
	PaginationValidator interfaces.PaginationValidator
	Repo                interfaces.FriendOfFriendPagingRepository
}

func NewFriendOfFriendPagingHandler(v interfaces.Validator, p interfaces.PaginationValidator, r interfaces.FriendOfFriendPagingRepository) *FriendOfFriendPagingHandler {
	return &FriendOfFriendPagingHandler{
		Validator:           v,
		PaginationValidator: p,
		Repo:                r,
	}
}

func (h *FriendOfFriendPagingHandler) FriendOfFriendPaging(c echo.Context) error {
	userID := c.QueryParam("id")
	if userID == "" || len(userID) > 20 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id must be a non-empty string up to 20 characters"})
	}

	exists, err := h.Validator.UserExists(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if !exists {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user not found"})
	}

	limit, page, err := h.PaginationValidator.ParseAndValidatePagination(c)
	if err != nil {
		if httpErr, ok := err.(*echo.HTTPError); ok {
			return c.JSON(httpErr.Code, map[string]string{"error": httpErr.Message.(string)})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	offset := (page - 1) * limit

	result, err := h.Repo.GetFriendOfFriendByIDWithPaging(userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(result) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no friends of friends found"})
	}

	return c.JSON(http.StatusOK, result)
}
