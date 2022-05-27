package main

import (
	"github.com/YJ9938/DouYin/model"
	"github.com/YJ9938/DouYin/router"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	model.InitMySQL()

	router.InitRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
