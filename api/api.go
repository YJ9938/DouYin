package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Common status code and status message.
type Status struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_msg"`
}

// A short hand to respond with error message.
func Error(c *gin.Context, statusCode int, statusMessage string) {
	c.JSON(http.StatusOK, Status{
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
	})
}

// A short hand to respond with internal error.
func InternalError(c *gin.Context) {
	Error(c, 500, "系统内部错误")
}
