package get

import (
	"net/http"
	"strconv"

	"minimal_sns_app/interfaces"

	"github.com/labstack/echo/v4"
)

type FriendOfFriendPagingHandler struct {
	UserValidator       interfaces.UserValidator
	PaginationValidator interfaces.PaginationValidator
	Repo                interfaces.FriendOfFriendPagingRepository
}

func NewFriendOfFriendPagingHandler(u interfaces.UserValidator, p interfaces.PaginationValidator, r interfaces.FriendOfFriendPagingRepository) *FriendOfFriendPagingHandler {
	return &FriendOfFriendPagingHandler{
		UserValidator:       u,
		PaginationValidator: p,
		Repo:                r,
	}
}

func (h *FriendOfFriendPagingHandler) FriendOfFriendPaging(c echo.Context) error {
	idStr := c.QueryParam("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id must be a positive integer"})
	}

	limit, page, err := h.PaginationValidator.ParseAndValidatePagination(c)
	if err != nil {
		return err
	}

	exist, err := h.UserValidator.UserExists(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if !exist {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user not found"})
	}

	offset := (page - 1) * limit
	result, err := h.Repo.GetFriendOfFriendByIDWithPaging(id, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(result) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "no friends of friends found"})
	}

	return c.JSON(http.StatusOK, result)
}
