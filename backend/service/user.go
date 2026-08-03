package service

import (
	"github.com/Anyeling0620/OnlineOJ/backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetUserDetail
// @Tags 公共方法
// @Summary 用户详情
// @Param identity query string false "用户唯一id"
// @Success 200 {object} SuccessResponse "成功返回列表数据"
// @Failure 400 {object} FailResponse "请求参数错误"
// @Failure 401 {object} FailResponse "未授权"
// @Failure 403 {object} FailResponse "权限不足"
// @Failure 404 {object} FailResponse "资源不存在"
// @Failure 500 {object} FailResponse "服务器内部错误"
// @Router /users/details [get] {}
func GetUserDetail(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "用户唯一标识不能为空",
		})
		return
	}
	user := &models.UserBasic{}
	err := models.DB.Omit("password").Where("identity=?", identity).First(&user).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "not found user identity:" + identity + " error:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    user,
		"message": "success",
	})
}
