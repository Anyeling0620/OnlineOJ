package service

import (
	"fmt"
	"github.com/Anyeling0620/OnlineOJ/backend/define"
	"github.com/Anyeling0620/OnlineOJ/backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	Code    int         `json:"code" example:"-1"`
	Data    interface{} `json:"data"`
	Message string      `json:"message" example:"fail"`
}

// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "请输入页数，默认为1"
// @Param size query int false "请输入每页结果个数，默认为20"
// @Param keyword query string false "模糊搜索关键词"
// @Success 200 {object} SuccessResponse "成功返回列表数据"
// @Failure 400 {object} FailResponse "请求参数错误"
// @Failure 401 {object} FailResponse "未授权"
// @Failure 403 {object} FailResponse "权限不足"
// @Failure 404 {object} FailResponse "资源不存在"
// @Failure 500 {object} FailResponse "服务器内部错误"
// @Router /problems [get] {}
func GetProblemList(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		fmt.Println("page conv error:", err)
		return
	}
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		fmt.Println("size conv error:", err)
		return
	}
	keyword := c.DefaultQuery("keyword", "")
	tx := models.GetProblemList(keyword)
	list := make([]*models.ProblemBasic, 0)
	offset := (page - 1) * size // 传入的参数从1开始，但是数据库的offset是从0开始

	var count int64
	err = tx.Count(&count).Omit("content").Offset(offset).Limit(size).Find(&list).Error

	if err != nil {
		fmt.Println("get problem list error:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"count": count,
			"list":  list,
		},
		"message": "success",
	})
}
