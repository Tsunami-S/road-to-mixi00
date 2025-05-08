package main

import (
	"minimal_sns_app/configs"
	"minimal_sns_app/db"
	"minimal_sns_app/handler/create"
	"minimal_sns_app/handler/get"
	"minimal_sns_app/handler/get_all"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func main() {
	db.InitDB()
	conf := configs.Get()

	e := echo.New()

	// ex00
	e.GET("/get_friend_list", get.Friend)
	// ex01,02
	e.GET("/get_friend_of_friend_list", get.FriendOfFriend)
	// ex03
	e.GET("/get_friend_of_friend_list_paging", get.FriendOfFriendPaging)

	// bonus
	e.POST("/add_new_user", create.AddNewUser)
	e.POST("/block_user", create.AddBlockList)
	e.POST("/request_friend", create.RequestFriend)
	e.POST("/respond_friend_request", create.RespondRequest)
	e.GET("/pending_requests", get.PendingRequests)

	// for debug
	e.GET("/get_all_users", get_all.Users)
	e.GET("/get_all_friends", get_all.FriendLinks)
	e.GET("/get_all_blocks", get_all.BlockList)
	e.GET("/get_all_requests", get_all.RequestList)

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
