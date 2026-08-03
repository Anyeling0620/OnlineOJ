package service

import (
	"fmt"
	"github.com/Anyeling0620/OnlineOJ/backend/define"
	"github.com/Anyeling0620/OnlineOJ/backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetSubmitList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "请输入页数，默认为1"
// @Param size query int false "请输入每页结果个数，默认为20"
// @Param problem_identity query string false "问题唯一id"
// @Param user_identity query string false "用户唯一id"
// @Param status query int false "状态"
// @Success 200 {object} SuccessResponse "成功返回列表数据"
// @Failure 400 {object} FailResponse "请求参数错误"
// @Failure 401 {object} FailResponse "未授权"
// @Failure 403 {object} FailResponse "权限不足"
// @Failure 404 {object} FailResponse "资源不存在"
// @Failure 500 {object} FailResponse "服务器内部错误"
// @Router /submit/lists [get] {}
func GetSubmitList(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		fmt.Println("page conv error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "page convert error",
		})
		return
	}
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		fmt.Println("size conv error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "size convert error",
		})
		return
	}
	problemIdentity := c.DefaultQuery("problem_identity", "")
	if problemIdentity == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "problemIdentity参数不能为空",
		})
		return
	}

	userIdentity := c.DefaultQuery("user_identity", "")
	if userIdentity == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "userIdentity参数不能为空",
		})
		return
	}

	status, err := strconv.Atoi(c.DefaultQuery("status", "-1"))
	if err != nil {
		fmt.Println("status conv error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "status conv error",
		})
		return
	}

	tx := models.GetSubmitList(problemIdentity, userIdentity, status)

	var count int64
	offset := (page - 1) * size
	list := make([]*models.SubmitBasic, 0)
	err = tx.Count(&count).Offset(offset).Limit(size).Find(&list).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "没找到对应的提交记录 err:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    gin.H{"count": count, "list": list},
		"message": "success",
	})
}
