package config

import (
	"tiktok/controller"
	"tiktok/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	// 主路由组
	g := r.Group("/douyin")
	{
		// user路由组
		userGroup := g.Group("/user")
		{
			userGroup.GET("/", middleware.JwtMiddleware(), controller.UserInfo)
			userGroup.POST("/login/", controller.UserLogin)
			userGroup.POST("/register/", controller.UserRegister)
		}

		// publish路由组
		publishGroup := g.Group("/publish")
		{
			publishGroup.POST("/action/", middleware.JwtMiddleware(), controller.Publish)
			publishGroup.GET("/list/", middleware.JwtMiddleware(), controller.PublishList)

		}

		// feed
		g.GET("/feed/", controller.Feed)

		favoriteGroup := g.Group("/favorite")
		{
			favoriteGroup.POST("/action/", middleware.JwtMiddleware(), controller.Like)
			favoriteGroup.GET("/list/", middleware.JwtMiddleware(), controller.LikeList)
		}

		// comment路由组
		commentGroup := g.Group("/comment")
		{
			commentGroup.POST("/action/", middleware.JwtMiddleware(), controller.Comment)
			commentGroup.GET("/list/", controller.CommentList)
		}
		followGroup := g.Group("/relation")
		{
			followGroup.POST("/action/", middleware.JwtMiddleware(), controller.Follow)
			followGroup.GET("/follow/list/", middleware.JwtMiddleware(), controller.FollowList)
			followGroup.GET("/follower/list/", middleware.JwtMiddleware(), controller.FollowerList)
			followGroup.GET("/friend/list/", middleware.JwtMiddleware(), controller.FriendList)
		}

		// message路由组
		messageGroup := g.Group("/message")
		{
			messageGroup.POST("/action/", middleware.JwtMiddleware(), controller.Message)
			messageGroup.GET("/chat/", middleware.JwtMiddleware(), controller.MessageList)
		}
	}

	return r
}
