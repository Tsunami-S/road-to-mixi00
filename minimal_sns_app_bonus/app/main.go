package main

import (
	"github.com/labstack/echo/v4"
	"minimal_sns_app/configs"
	"minimal_sns_app/db"
	"minimal_sns_app/handler"
	"net/http"
	"strconv"
)

func main() {
	db.InitDB()
	conf := configs.Get()

	e := echo.New()

	// ex00
	e.GET("/get_friend_list", handler.GetFriendList)
	// ex01,02
	e.GET("/get_friend_of_friend_list", handler.GetFriendOfFriendList)
	// ex03
	e.GET("/get_friend_of_friend_list_paging", handler.GetFriendOfFriendListPaging)

	// bonus
	e.GET("/add_new_user", handler.AddNewUser)
	e.GET("/block_user", handler.AddBlockList)
	e.GET("/request_friend", handler.RequestFriend)
	e.GET("/respond_friend_request", handler.RespondFriendRequest)
	e.GET("/pending_requests", handler.GetPendingRequests)

	// for debug
	e.GET("/all_users", handler.GetAllUsers)
	e.GET("/all_friends", handler.GetAllFriendLinks)
	e.GET("/all_blocks", handler.GetAllBlockList)
	e.GET("/all_requests", handler.GetAllFriendRequests)

	// for error
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}
		code := http.StatusNotFound
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		c.JSON(code, map[string]string{"error": "invalid endpoint"})
	}

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(conf.Server.Port)))
}
