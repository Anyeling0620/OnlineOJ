package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Anyeling0620/OnlineOJ/backend/define"
	"github.com/Anyeling0620/OnlineOJ/backend/models"
	"github.com/Anyeling0620/OnlineOJ/backend/util"
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
// @Router /problems [get] {}
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
	page = max(page, 0)
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
// @Param problemIdentity path string true "问题唯一标识"
// @Success 200 {object} SuccessResponse "成功返回问题详情"
// @Failure 400 {object} FailResponse "请求参数错误"
// @Failure 401 {object} FailResponse "未授权"
// @Failure 403 {object} FailResponse "权限不足"
// @Failure 404 {object} FailResponse "资源不存在"
// @Failure 500 {object} FailResponse "服务器内部错误"
// @Router /problems/{problemIdentity} [get]
func GetProblemDetail(c *gin.Context) {
	identity := c.Param("problemIdentity")
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

// ProblemCreate
// @Tags 管理员私有方法
// @Summary 问题创建
// @Param Authorization header string true "Authorization"
// @Param title formData string true "标题"
// @Param content formData string true "内容"
// @Param max_runtime formData string true "最大运行时间(ms)"
// @Param max_mem formData string true "最大运行内存(kb)"
// @Param category_ids formData []int false "分类ID数组" collectionFormat(multi)
// @Param test_cases formData []string true "测试样例数组" collectionFormat(multi)
// @Success 200 {object} SuccessResponse "成功返回问题详情"
// @Failure 400 {object} FailResponse "请求参数错误"
// @Failure 401 {object} FailResponse "未授权"
// @Failure 403 {object} FailResponse "权限不足"
// @Failure 404 {object} FailResponse "资源不存在"
// @Failure 500 {object} FailResponse "服务器内部错误"
// @Router /problems [post] {}
func ProblemCreate(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	maxRuntime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	maxMem, _ := strconv.Atoi(c.PostForm("max_mem"))
	categoryIds := c.PostFormArray("category_ids")
	testCases := c.PostFormArray("test_cases")
	if title == "" || content == "" || len(categoryIds) <= 0 || len(testCases) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数不能为空",
		})
		return
	}
	identity := util.GetUUID()
	data := &models.ProblemBasic{
		Title:      title,
		Identity:   identity,
		Content:    content,
		MaxRuntime: maxRuntime,
		MaxMem:     maxMem,
	}
	// 创建问题的分类
	categoryBasics := make([]*models.ProblemCategory, 0)
	for _, categoryId := range categoryIds {
		intCategoryId, err := strconv.Atoi(categoryId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "category id convert error",
			})
			return
		}
		if intCategoryId <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "分类ID不能为负数",
			})
		}
		categoryBasics = append(categoryBasics, &models.ProblemCategory{
			ProblemID:  data.ID,
			CategoryID: uint(intCategoryId),
		})
	}
	data.ProblemCategories = categoryBasics
	// 创建问题的测试用例
	testCasesBasics := make([]*models.TestCaseBasic, 0)
	for _, testCase := range testCases {
		// case : {"input":"1 2\n", "output": "3\n"}
		caseMap := make(map[string]string)
		err := json.Unmarshal([]byte(testCase), &caseMap)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "测试样例格式错误",
			})
			return
		}
		if _, ok := caseMap["input"]; !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "测试样例输入格式错误",
			})
			return
		}
		if _, ok := caseMap["output"]; !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "测试样例输出格式错误",
			})
			return
		}

		testCaseBasic := &models.TestCaseBasic{
			Identity:        util.GetUUID(),
			ProblemIdentity: identity,
			Input:           caseMap["input"],
			Output:          caseMap["output"],
		}
		testCasesBasics = append(testCasesBasics, testCaseBasic)
	}
	data.TestCaseBasics = testCasesBasics

	err := models.DB.Create(data).Error
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": "问题创建失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"list": data,
		},
		"message": "success",
	})
}

// ProblemModify
// @Tags 管理员私有方法
// @Summary 问题修改
// @Param Authorization header string true "Authorization"
// @Param problemIdentity path string true "问题唯一标识"
// @Param title formData string true "标题"
// @Param content formData string true "内容"
// @Param max_runtime formData string true "最大运行时间(ms)"
// @Param max_mem formData string true "最大运行内存(kb)"
// @Param category_ids formData []int false "分类ID数组" collectionFormat(multi)
// @Param test_cases formData []string true "测试样例数组" collectionFormat(multi)
// @Success 200 {object} SuccessResponse "成功返回问题详情"
// @Failure 400 {object} FailResponse "请求参数错误"
// @Failure 401 {object} FailResponse "未授权"
// @Failure 403 {object} FailResponse "权限不足"
// @Failure 404 {object} FailResponse "资源不存在"
// @Failure 500 {object} FailResponse "服务器内部错误"
// @Router /problems/{problemIdentity} [put]
func ProblemModify(c *gin.Context) {
	identity := c.Param("problemIdentity")
	title := c.PostForm("title")
	content := c.PostForm("content")
	maxRuntime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	maxMem, _ := strconv.Atoi(c.PostForm("max_mem"))
	categoryIds := c.PostFormArray("category_ids")
	testCases := c.PostFormArray("test_cases")
	if title == "" || content == "" || len(categoryIds) <= 0 || len(testCases) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数不能为空",
		})
		return
	}

	if err := models.DB.Transaction(func(tx *gorm.DB) error {
		// 要么全成功 要么全不成功
		// 问题基本信息保存
		problemBasic := &models.ProblemBasic{
			Identity:   identity,
			Title:      title,
			Content:    content,
			MaxRuntime: maxRuntime,
			MaxMem:     maxMem,
		}
		err := tx.Where("identity=?", identity).Updates(problemBasic).Error
		if err != nil {
			fmt.Println(err)
			return err
		}
		//    查询问题详情
		err = tx.Where("identity=?", identity).Find(problemBasic).Error
		if err != nil {
			fmt.Println(err)
			return err
		}
		// 关联问题分类更新
		// 1. 删除已存在的关联关系
		err = tx.Where("problem_id=?", problemBasic.ID).Delete(&models.ProblemCategory{}).Error
		if err != nil {
			fmt.Println(err)
			return err
		}
		// 2. 新增新的关联关系
		pcs := make([]*models.ProblemCategory, 0)
		for _, id := range categoryIds {
			intId, err := strconv.Atoi(id)
			if err != nil {
				fmt.Println(err)
				return err
			}
			pcs = append(pcs, &models.ProblemCategory{
				ProblemID:  problemBasic.ID,
				CategoryID: uint(intId),
			})
		}
		err = tx.Create(&pcs).Error
		if err != nil {
			fmt.Println(err)
			return err
		}
		// 关联测试样例更新
		// 1. 删除已存在的关联关系
		err = tx.Where("problem_identity=?", identity).Delete(&models.TestCaseBasic{}).Error
		if err != nil {
			fmt.Println(err)
			return err
		}
		// 2. 增加新的关联关系
		tcs := make([]*models.TestCaseBasic, 0)
		for _, testCase := range testCases {
			caseMap := make(map[string]string)
			err = json.Unmarshal([]byte(testCase), &caseMap)
			if err != nil {
				fmt.Println(err)
				return err
			}
			if _, ok := caseMap["input"]; !ok {
				return errors.New("测试案例input格式错误")
			}
			if _, ok := caseMap["output"]; !ok {
				return errors.New("测试案例output格式错误")
			}
			tcs = append(tcs, &models.TestCaseBasic{
				Identity:        util.GetUUID(),
				ProblemIdentity: identity,
				Input:           caseMap["input"],
				Output:          caseMap["output"],
			})
		}
		err = tx.Create(&tcs).Error
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadRequest,
			"message": "problem modify error" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
	})

}
