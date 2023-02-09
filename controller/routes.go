package controller

import (
	middleware "tiktok/pkg/mw"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	// 主路由组
	g := r.Group("/")
	{
		// user路由组
		userGroup := g.Group("/user")
		{
			userGroup.GET("/", middleware.JwtMiddleware(), UserInfo)
			userGroup.POST("/login/", UserLogin)
			userGroup.POST("/register/", UserRegister)
		}

		// publish路由组
		// publishGroup := g.Group("/publish")
		// {
		// 	publishGroup.POST("/action/", middleware.JwtMiddleware(), Publish)
		// 	publishGroup.GET("/list/", middleware.JwtMiddleware(), controller.PublishList)

		// }

		// // feed
		// g.GET("/feed/", controller.Feed)

		// favoriteGroup := g.Group("favorite")
		// {
		// 	favoriteGroup.POST("/action/", middleware.JwtMiddleware(), controller.Favorite)
		// 	favoriteGroup.GET("/list/", middleware.JwtMiddleware(), controller.FavoriteList)
		// }

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
