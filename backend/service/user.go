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
	"time"
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

// SendCode
// @Tags 公共方法
// @Summary 验证码发送
// @Param mail formData string true "mail"
// @Success 200 {object} SuccessResponse "成功返回列表数据"
// @Failure 400 {object} FailResponse "请求参数错误"
// @Failure 401 {object} FailResponse "未授权"
// @Failure 403 {object} FailResponse "权限不足"
// @Failure 404 {object} FailResponse "资源不存在"
// @Failure 500 {object} FailResponse "服务器内部错误"
// @Router /users/code [post] {}
func SendCode(c *gin.Context) {
	mail := c.PostForm("mail")
	if mail == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "mail is null",
		})
		return
	}

	code := util.GetRand()
	models.RDB.Set(c, mail, code, time.Second*60)
	err := util.SendCode(mail, code)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": "send code error:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
	})
}

// Register
// @Tags 公共方法
// @Summary 用户注册
// @Param name formData string true "name"
// @Param password formData string true "password"
// @Param phone formData string false "phone"
// @Param mail formData string true "mail"
// @Param code formData string true "code"
// @Success 200 {object} SuccessResponse "成功返回列表数据"
// @Failure 400 {object} FailResponse "请求参数错误"
// @Failure 401 {object} FailResponse "未授权"
// @Failure 403 {object} FailResponse "权限不足"
// @Failure 404 {object} FailResponse "资源不存在"
// @Failure 500 {object} FailResponse "服务器内部错误"
// @Router /users/register [post] {}
func Register(c *gin.Context) {
	mail := c.PostForm("mail")
	userCode := c.PostForm("code")
	if mail == "" || userCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "验证码不正确",
		})
		return
	}

	name := c.PostForm("name")
	password := c.PostForm("password")
	if name == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "name or password is null",
		})
		return
	}
	phone := c.PostForm("phone")

	sysCode, err := models.RDB.Get(c, mail).Result()
	if err != nil {
		log.Println("Get Code Error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "验证码不正确，请重新获取验证码",
		})
		return
	}
	// 判断邮箱是否存在
	var cnt int64
	err = models.DB.Where("mail=?", mail).Model(&models.UserBasic{}).Count(&cnt).Error
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": "get userinfo error",
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "该邮箱已被注册",
		})
		return
	}

	err = models.DB.Where("name=?", name).Model(&models.UserBasic{}).Count(&cnt).Error
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": "该用户名已被注册",
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "this name is exist",
		})
		return
	}

	if sysCode != userCode {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "验证码不正确",
		})
		return
	}

	userIdentity := util.GetUUID()
	data := &models.UserBasic{
		Identity: userIdentity,
		Name:     name,
		Password: util.GetMd5(password),
		Phone:    phone,
		Mail:     mail,
	}
	err = models.DB.Create(&data).Error
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": "创建用户时出错",
		})
		return
	}

	// 生成token
	token, err := util.GenerateToken(userIdentity, name)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": "generate token error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"token": token,
		},
		"message": "success",
	})
}
