package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

func Error(c *gin.Context, statusCode int32, statusMsg string) {
	c.JSON(http.StatusOK, Response{
		StatusCode: statusCode,
		StatusMsg:  statusMsg,
	})
}

// A short hand to respond with internal error.
func InternalError(c *gin.Context) {
	Error(c, 500, "系统内部错误")
}
