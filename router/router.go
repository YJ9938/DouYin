package router

import (
	"github.com/YJ9938/DouYin/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")

	// Unauthorized APIs
	{
		// public directory is used to serve static resources
		r.Static("/static", "./public")

		// user module
		apiRouter.GET("/feed/", controller.Feed)
		apiRouter.POST("/user/register/", controller.Register)
		apiRouter.POST("/user/login/", controller.Login)
	}

	// Authorized APIs without token query parameter
	apiRouter.POST("/publish/action/", controller.Publish)

	// Authorized APIs with token query parameter
	userRouter := r.Group("/douyin/user")
	userRouter.Use(controller.AuthMiddleware)
	{
		userRouter.GET("/", controller.UserInfo)
	}

	publishGroup := r.Group("/douyin/publish")
	publishGroup.Use(controller.AuthMiddleware)
	{
		publishGroup.GET("/list/", controller.PublishList)
	}

	favoriteGroup := r.Group("/douyin/favorite")
	favoriteGroup.Use(controller.AuthMiddleware)
	{
		favoriteGroup.POST("/action/", controller.FavoriteAction)
		favoriteGroup.GET("/list/", controller.FavoriteList)
	}

	commentGroup := r.Group("/douyin/comment")
	commentGroup.Use(controller.AuthMiddleware)
	{
		commentGroup.POST("/action/", controller.CommentAction)
		commentGroup.GET("/list/", controller.CommentList)
	}

	relationGroup := r.Group("/douyin/relation")
	relationGroup.Use(controller.AuthMiddleware)
	{
		relationGroup.POST("/action/", controller.RelationAction)
		relationGroup.GET("/follow/list/", controller.FollowList)
		relationGroup.GET("/follower/list/", controller.FollowerList)
	}
}
