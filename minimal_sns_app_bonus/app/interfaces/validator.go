package interfaces

import "github.com/labstack/echo/v4"

type Validator interface {
	UserExists(id string) (bool, error)
}

type PaginationValidator interface {
	ParseAndValidatePagination(c echo.Context) (limit int, page int, err error)
}
