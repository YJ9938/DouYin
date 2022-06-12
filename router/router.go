package router

import (
	"github.com/YJ9938/DouYin/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	// Unauthorized APIs
	noAuthAPIs := r.Group("/douyin")
	{
		// public directory is used to serve static resources
		r.Static("/static", "./public")

		// user module
		noAuthAPIs.GET("/feed/", controller.Feed)
		noAuthAPIs.POST("/user/register/", controller.Register)
		noAuthAPIs.POST("/user/login/", controller.Login)
	}

	// Form auth APIs
	formAuthAPIs := r.Group("/douyin")
	formAuthAPIs.Use(controller.FormAuthMiddleware)
	{
		formAuthAPIs.POST("/publish/action/", controller.Publish)
	}

	// Query auth APIs
	queryAuthAPIs := r.Group("/douyin")
	queryAuthAPIs.Use(controller.QueryAuthMiddleware)
	{
		// user module
		queryAuthAPIs.GET("/user/", controller.UserInfo)
		// publish module
		queryAuthAPIs.GET("/publish/list/", controller.PublishList)
		// favorite module
		queryAuthAPIs.POST("/favorite/action/", controller.FavoriteAction)
		queryAuthAPIs.GET("/favorite/list/", controller.FavoriteList)
		// comment module
		queryAuthAPIs.POST("/comment/action/", controller.CommentAction)
		queryAuthAPIs.GET("/comment/list/", controller.CommentList)
		// relation module
		queryAuthAPIs.POST("/relation/action/", controller.RelationAction)
		queryAuthAPIs.GET("/relation/follow/list/", controller.FollowList)
		queryAuthAPIs.GET("/relation/follower/list/", controller.FollowerList)
	}
}
