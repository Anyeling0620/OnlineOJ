package service

import (
	"errors"
	"fmt"
	"github.com/Anyeling0620/OnlineOJ/backend/define"
	"github.com/Anyeling0620/OnlineOJ/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "请输入页数，默认为1"
// @Param size query int false "请输入每页结果个数，默认为20"
// @Param keyword query string false "模糊搜索关键词"
// @Param category_identity query string false "分类唯一id"
// @Success 200 {object} SuccessResponse "成功返回列表数据"
// @Failure 400 {object} FailResponse "请求参数错误"
// @Failure 401 {object} FailResponse "未授权"
// @Failure 403 {object} FailResponse "权限不足"
// @Failure 404 {object} FailResponse "资源不存在"
// @Failure 500 {object} FailResponse "服务器内部错误"
// @Router /problems/lists [get] {}
func GetProblemList(c *gin.Context) {
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
	keyword := c.DefaultQuery("keyword", "")

	categoryIdentity := c.DefaultQuery("category_identity", "")

	tx := models.GetProblemList(keyword, categoryIdentity)
	list := make([]*models.ProblemBasic, 0)
	offset := (page - 1) * size // 传入的参数从1开始，但是数据库的offset是从0开始

	var count int64
	err = tx.Count(&count).Omit("content").Offset(offset).Limit(size).Find(&list).Error

	if err != nil {
		fmt.Println("get problem list from db error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "page convert error",
		})
		return
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

// GetProblemDetail
// @Tags 公共方法
// @Summary 问题详情
// @Param problem_identity query string false "问题唯一id"
// @Success 200 {object} SuccessResponse "成功返回问题详情"
// @Failure 400 {object} FailResponse "请求参数错误"
// @Failure 401 {object} FailResponse "未授权"
// @Failure 403 {object} FailResponse "权限不足"
// @Failure 404 {object} FailResponse "资源不存在"
// @Failure 500 {object} FailResponse "服务器内部错误"
// @Router /problems/details [get] {}
func GetProblemDetail(c *gin.Context) {
	identity := c.Query("problem_identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "问题唯一标识不能为空",
		})
		return
	}
	data := new(models.ProblemBasic)
	err := models.DB.Where("identity = ?", identity).Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").First(data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusNotFound,
				"message": "问题不存在",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "非法的uuid",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    data,
		"message": "success",
	})
}
