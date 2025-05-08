package main

import (
	"minimal_sns_app/configs"
	"minimal_sns_app/db"

	"minimal_sns_app/handler/create"
	"minimal_sns_app/handler/get"
	"minimal_sns_app/handler/get_all"

	repo_create "minimal_sns_app/repository/create"
	repo_get "minimal_sns_app/repository/get"

	"minimal_sns_app/handler/validate"

	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func main() {
	db.InitDB()
	conf := configs.Get()
	e := echo.New()

	validator := &validate.RealValidator{}
	paginationValidator := &validate.RealPaginationValidator{}

	friendHandler := get.NewFriendHandler(validator, &repo_get.RealFriendRepository{})
	fofHandler := get.NewFriendOfFriendHandler(validator, &repo_get.RealFriendOfFriendRepository{})
	fofPagingHandler := get.NewFriendOfFriendPagingHandler(validator, paginationValidator, &repo_get.RealFriendOfFriendPagingRepository{})
	pendingHandler := get.NewPendingRequestHandler(validator, &repo_get.RealFriendRequestRepository{})

	addUserHandler := create.NewUserHandler(&repo_create.RealUserRepository{})
	blockHandler := create.NewBlockHandler(validator, &repo_create.RealBlockRepository{})
	requestHandler := create.NewRequestFriendHandler(validator, &repo_create.RealFriendRequestRepository{})
	respondHandler := create.NewFriendRespondHandler(
		validator,
		&repo_create.RealFriendRespondRepository{},
	)

	// ex00
	e.GET("/get_friend_list", friendHandler.Friend)

	// ex01, 02
	e.GET("/get_friend_of_friend_list", fofHandler.FriendOfFriend)

	// ex03
	e.GET("/get_friend_of_friend_list_paging", fofPagingHandler.FriendOfFriendPaging)

	// bonus
	e.POST("/add_new_user", addUserHandler.AddNewUser)
	e.POST("/block_user", blockHandler.AddBlockList)
	e.POST("/request_friend", requestHandler.RequestFriend)
	e.POST("/respond_friend_request", respondHandler.RespondRequest)
	e.GET("/pending_requests", pendingHandler.PendingRequests)

	// debug routes
	e.GET("/get_all_users", get_all.Users)
	e.GET("/get_all_friends", get_all.FriendLinks)
	e.GET("/get_all_blocks", get_all.BlockList)
	e.GET("/get_all_requests", get_all.RequestList)

	// fallback error handler
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
