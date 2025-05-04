package handler

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/db"
	"minimal_sns_app/model"
	"net/http"
)

func GetAllUsers(c echo.Context) error {
	var users []model.User
	if err := db.DB.Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch users"})
	}
	return c.JSON(http.StatusOK, users)
}

func GetAllFriendLinks(c echo.Context) error {
	var links []model.FriendLink
	if err := db.DB.Find(&links).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch friend links"})
	}
	return c.JSON(http.StatusOK, links)
}

func GetAllBlockList(c echo.Context) error {
	var blocks []model.BlockList
	if err := db.DB.Find(&blocks).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch block list"})
	}
	return c.JSON(http.StatusOK, blocks)
}

func GetAllFriendRequests(c echo.Context) error {
	var requests []model.FriendRequest
	if err := db.DB.Find(&requests).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch friend requests"})
	}
	return c.JSON(http.StatusOK, requests)
}
