package service

import (
	"errors"
	"fmt"
	"github.com/Anyeling0620/OnlineOJ/backend/models"
	"github.com/Anyeling0620/OnlineOJ/backend/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
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

// Login
// @Tags 公共方法
// @Summary 用户登录
// @Param username formData string false "用户昵称"
// @Param password formData string false "用户密码"
// @Success 200 {object} SuccessResponse "成功返回列表数据"
// @Failure 400 {object} FailResponse "请求参数错误"
// @Failure 401 {object} FailResponse "未授权"
// @Failure 403 {object} FailResponse "权限不足"
// @Failure 404 {object} FailResponse "资源不存在"
// @Failure 500 {object} FailResponse "服务器内部错误"
// @Router /users/login [post] {}
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "username or password is null",
		})
	}
	password = util.GetMd5(password)
	log.Print("username: ", username, " password :", password)

	data := new(models.UserBasic)
	err := models.DB.Where("name=? AND password=?", username, password).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("用户名或密码错误")
			c.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": "用户名或密码错误",
			})
			return
		}

		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": "get userinfo error:" + err.Error(),
		})
		return
	}

	token, err := util.GenerateToken(data.Identity, data.Name)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": "generate token error:" + err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"token": token,
		},
		"message": "success",
	})
}
