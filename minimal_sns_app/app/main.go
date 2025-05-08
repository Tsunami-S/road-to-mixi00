package main

import (
	"net/http"
	"strconv"

	"minimal_sns_app/configs"
	"minimal_sns_app/db"
	"minimal_sns_app/handler/get"
	"minimal_sns_app/handler/get_all"
	repo_get "minimal_sns_app/repository/get"
	"minimal_sns_app/repository/validate"

	"github.com/labstack/echo/v4"
)

func main() {
	// DB初期化
	db.InitDB()

	// 設定取得
	conf := configs.Get()
	e := echo.New()

	// ✅ ハンドラの依存注入と初期化

	friendHandler := get.NewFriendHandler(
		&validate.RealValidator{},
		&repo_get.RealFriendRepository{},
	)

	friendOfFriendHandler := get.NewFriendOfFriendHandler(
		&validate.RealValidator{},
		&repo_get.RealFriendOfFriendRepository{},
	)

	friendOfFriendPagingHandler := get.NewFriendOfFriendPagingHandler(
		&validate.RealValidator{},
		&validate.RealPaginationValidator{},
		&repo_get.RealFriendOfFriendPagingRepository{},
	)

	// ✅ ルーティング
	e.GET("/get_friend_list", friendHandler.Friend)
	e.GET("/get_friend_of_friend_list", friendOfFriendHandler.FriendOfFriend)
	e.GET("/get_friend_of_friend_list_paging", friendOfFriendPagingHandler.FriendOfFriendPaging)

	// ✅ デバッグ用エンドポイント
	e.GET("/get_all_users", get_all.Users)
	e.GET("/get_all_friends", get_all.FriendLinks)
	e.GET("/get_all_blocks", get_all.BlockList)

	// ✅ 無効エンドポイント処理
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

	// ✅ 起動
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(conf.Server.Port)))
}
