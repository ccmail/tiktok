package controller

import (
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
			userGroup.GET("/", middleware.JwtMiddleware(), UserInfo)
			userGroup.POST("/login/", UserLogin)
			userGroup.POST("/register/", UserRegister)
		}

		// publish路由组
		publishGroup := g.Group("/publish")
		{
			publishGroup.POST("/action/", middleware.JwtMiddleware(), Publish)
			publishGroup.GET("/list/", middleware.JwtMiddleware(), PublishList)

		}

		// feed
		g.GET("/feed/", Feed)

		favoriteGroup := g.Group("/favorite")
		{
			favoriteGroup.POST("/action/", middleware.JwtMiddleware(), Like)
			favoriteGroup.GET("/list/", middleware.JwtMiddleware(), LikeList)
		}

		// comment路由组
		commentGroup := g.Group("/comment")
		{
			commentGroup.POST("/action/", middleware.JwtMiddleware(), Comment)
			commentGroup.GET("/list/", CommentList)
		}
		followGroup := g.Group("/relation")
		{
			followGroup.POST("/action/", middleware.JwtMiddleware(), Follow)
			followGroup.GET("/follow/list/", middleware.JwtMiddleware(), FollowList)
			followGroup.GET("/follower/list/", middleware.JwtMiddleware(), FollowerList)
			followGroup.GET("/friend/list/", middleware.JwtMiddleware(), FriendList)
		}

		// message路由组
		messageGroup := g.Group("/message")
		{
			messageGroup.POST("/action/", middleware.JwtMiddleware(), Message)
			messageGroup.GET("/chat/", middleware.JwtMiddleware(), MessageList)
		}
	}

	return r
}
