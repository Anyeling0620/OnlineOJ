package service

import (
	"fmt"
	"github.com/Anyeling0620/OnlineOJ/backend/define"
	"github.com/Anyeling0620/OnlineOJ/backend/models"
	"github.com/Anyeling0620/OnlineOJ/backend/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// GetCategoryList
// @Tags 管理员私有方法
// @Summary 分类列表
// @Param page query int false "请输入页数，默认为1"
// @Param size query int false "请输入每页结果个数，默认为20"
// @Param keyword query string false "模糊搜索关键词"
// @Param Authorization header string true "Authorization"
// @Success 200 {object} SuccessResponse "成功返回列表数据"
// @Failure 400 {object} FailResponse "请求参数错误"
// @Failure 401 {object} FailResponse "未授权"
// @Failure 403 {object} FailResponse "权限不足"
// @Failure 404 {object} FailResponse "资源不存在"
// @Failure 500 {object} FailResponse "服务器内部错误"
// @Router /categories [get] {}
func GetCategoryList(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		fmt.Println("page conv error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "page convert error",
		})
		return
	}
	page = max(page, 1)
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		fmt.Println("size conv error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "size convert error",
		})
		return
	}
	size = max(size, 1)
	keyword := c.DefaultQuery("keyword", "")

	var count int64
	categoryList := make([]*models.CategoryBasic, 0)
	err = models.DB.Model(new(models.CategoryBasic)).Omit("problem_categories").Where("name like ?", "%"+keyword+"%").Count(&count).Limit(size).Offset((page - 1) * size).Limit(size).Find(&categoryList).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "get category list error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"count": count,
			"list":  categoryList,
		},
		"message": "success",
	})
}

// CategoryCreate
// @Tags 管理员私有方法
// @Summary 分类创建
// @Param Authorization header string true "Authorization"
// @Param name formData string true "name"
// @Param parent_id formData string false "parent_id"
// @Success 200 {object} SuccessResponse "成功返回问题详情"
// @Failure 400 {object} FailResponse "请求参数错误"
// @Failure 401 {object} FailResponse "未授权"
// @Failure 403 {object} FailResponse "权限不足"
// @Failure 404 {object} FailResponse "资源不存在"
// @Failure 500 {object} FailResponse "服务器内部错误"
// @Router /categories [post] {}
func CategoryCreate(c *gin.Context) {
	name := c.PostForm("name")
	if len(name) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "name is null",
		})
		return
	}
	parentId, err := strconv.Atoi(c.DefaultPostForm("parent_id", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "parent_id conv fail",
		})
		return
	}
	identity := util.GetUUID()
	category := &models.CategoryBasic{
		Identity: identity,
		Name:     name,
		ParentId: parentId,
	}
	err = models.DB.Model(new(models.CategoryBasic)).Create(category).Error
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": "创建分类失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
	})
}

// CategoryModify
// @Tags 管理员私有方法
// @Summary 分类修改
// @Param Authorization header string true "Authorization"
// @Param categoryIdentity path string true "问题唯一标识"
// @Param name formData string true "name"
// @Param parent_id formData string false "parent_id"
// @Success 200 {object} SuccessResponse "成功返回问题详情"
// @Failure 400 {object} FailResponse "请求参数错误"
// @Failure 401 {object} FailResponse "未授权"
// @Failure 403 {object} FailResponse "权限不足"
// @Failure 404 {object} FailResponse "资源不存在"
// @Failure 500 {object} FailResponse "服务器内部错误"
// @Router /categories/{categoryIdentity} [put]
func CategoryModify(c *gin.Context) {
	identity := c.Param("categoryIdentity")
	if len(identity) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "identity is null",
		})
		return
	}
	name := c.PostForm("name")
	if len(name) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "name is null",
		})
		return
	}
	parentId, err := strconv.Atoi(c.DefaultPostForm("parent_id", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "parent_id conv fail",
		})
		return
	}
	category := &models.CategoryBasic{
		Identity: identity,
		Name:     name,
		ParentId: parentId,
	}
	err = models.DB.Model(new(models.CategoryBasic)).Where("identity=?", identity).Updates(category).Error
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": "分类修改失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
	})
}

// CategoryDelete
// @Tags 管理员私有方法
// @Summary 分类删除
// @Param Authorization header string true "Authorization"
// @Param categoryIdentity path string true "分类唯一标识"
// @Success 200 {object} SuccessResponse "成功返回问题详情"
// @Failure 400 {object} FailResponse "请求参数错误"
// @Failure 401 {object} FailResponse "未授权"
// @Failure 403 {object} FailResponse "权限不足"
// @Failure 404 {object} FailResponse "资源不存在"
// @Failure 500 {object} FailResponse "服务器内部错误"
// @Router /categories/{categoryIdentity} [delete]
func CategoryDelete(c *gin.Context) {
	identity := c.Param("categoryIdentity")
	if len(identity) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "identity is null",
		})
		return
	}

	var cnt int64
	err := models.DB.Model(new(models.ProblemCategory)).Where("category_id=(select id from category_basic cb where cb.identity = ? limit 1)", identity).Count(&cnt).Error
	if err != nil {
		log.Println("Get ProblemCategory fail:", err)
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": "获取分类关联的问题失败",
		})
		return
	}

	if cnt > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "该分类下还有问题，不可直接删除",
		})
		return
	}
	err = models.DB.Model(new(models.CategoryBasic)).Where("identity=?", identity).Delete(&models.CategoryBasic{}).Error
	if err != nil {
		log.Println("Delete ProblemCategory fail:", err)
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": "删除分类失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
	})
}
