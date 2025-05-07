package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/repository"
	"net/http"
)

func GetAllUsers(c echo.Context) error {
	users, err := repository.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch users"})
	}
	return c.JSON(http.StatusOK, users)
}

func GetAllFriendLinks(c echo.Context) error {
	links, err := repository.GetAllFriendLinks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch friend links"})
	}
	return c.JSON(http.StatusOK, links)
}

func GetAllBlockList(c echo.Context) error {
	blocks, err := repository.GetAllBlockList()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch block list"})
	}
	return c.JSON(http.StatusOK, blocks)
}

func GetAllFriendRequests(c echo.Context) error {
	requests, err := repository.GetAllFriendRequests()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch friend requests"})
	}
	return c.JSON(http.StatusOK, requests)
}
