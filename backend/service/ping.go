package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// SuccessResponse 成功响应
type SuccessResponse struct {
	Code    int      `json:"code" example:"200"`
	Data    ListData `json:"data"`
	Message string   `json:"message" example:"success"`
}

// ListData data内容
type ListData struct {
	Count int         `json:"count" example:"100"`
	List  interface{} `json:"list"`
}

// FailResponse 失败响应
type FailResponse struct {
	Code    int    `json:"code" example:"-1"`
	Message string `json:"message" example:"something went wrong"`
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":    http.StatusOK,
		"data":    "",
		"message": "pong",
	})
}
