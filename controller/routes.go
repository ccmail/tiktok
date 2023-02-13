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
			//publishGroup.GET("/list/", middleware.JwtMiddleware(), PublishList)
			publishGroup.GET("/list/", middleware.JwtMiddleware(), PublishList)

		}

		// feed
		g.GET("/feed/", Feed)

		favoriteGroup := g.Group("/favorite")
		{
			favoriteGroup.POST("/action/", middleware.JwtMiddleware(), Like)
			favoriteGroup.GET("/list/", middleware.JwtMiddleware(), LikeList)
		}

		followGroup := g.Group("/relation")
		{
			followGroup.POST("/action/", middleware.JwtMiddleware(), Follow)
			followGroup.GET("/follow/list/", middleware.JwtMiddleware(), FollowList)
			//followGroup.GET("/follow/list/", FollowList)
		}
		// // comment路由组
		// commentGroup := g.Group("/comment")
		// {
		// 	commentGroup.POST("/action/", middleware.JwtMiddleware(), controller.CommentAction)
		// 	commentGroup.GET("/list/", middleware.JwtMiddleware(), controller.CommentList)
		// }

		// // relation路由组
		// relationGroup := g.Group("relation")
		// {
		// 	relationGroup.POST("/action/", middleware.JwtMiddleware(), controller.RelationAction)
		// 	relationGroup.GET("/follow/list/", middleware.JwtMiddleware(), controller.FollowList)
		// 	relationGroup.GET("/follower/list/", middleware.JwtMiddleware(), controller.FollowerList)
		// }
	}

	return r
}
