package main

import (
	"github.com/YJ9938/DouYin/api"
	"github.com/YJ9938/DouYin/config"
	_ "github.com/YJ9938/DouYin/model"
	"github.com/gin-gonic/gin"
)

func main() {
	c := gin.Default()

	// Unauthorized APIs
	{
		c.POST("douyin/user/register/", api.Register)
		c.POST("douyin/user/login/", api.Login)
		c.Static("static/", config.C.Static.Path)
	}

	// Authorized APIs all contain a query parameter
	authGroup := c.Group("/douyin")
	authGroup.Use(api.AuthMiddleware)
	{
		authGroup.GET("/user/", api.QueryUserInfo)
	}

	c.Run()
}
