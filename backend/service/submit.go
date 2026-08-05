package service

import (
	"bytes"
	"fmt"
	"github.com/Anyeling0620/OnlineOJ/backend/define"
	"github.com/Anyeling0620/OnlineOJ/backend/models"
	"github.com/Anyeling0620/OnlineOJ/backend/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// GetSubmitList
// @Tags 公共方法
// @Summary 提交列表
// @Param page query int false "请输入页数，默认为1"
// @Param size query int false "请输入每页结果个数，默认为20"
// @Param problemIdentity query string true "问题唯一id"
// @Param userIdentity query string true "用户唯一id"
// @Param status query int false "状态"
// @Success 200 {object} SuccessResponse "成功返回列表数据"
// @Failure 400 {object} FailResponse "请求参数错误"
// @Failure 401 {object} FailResponse "未授权"
// @Failure 403 {object} FailResponse "权限不足"
// @Failure 404 {object} FailResponse "资源不存在"
// @Failure 500 {object} FailResponse "服务器内部错误"
// @Router /submissions [get] {}
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
	problemIdentity := c.DefaultQuery("problemIdentity", "")
	if problemIdentity == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "problemIdentity参数不能为空",
		})
		return
	}

	userIdentity := c.DefaultQuery("userIdentity", "")
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

// Submit
// @Tags 用户私有方法
// @Summary 代码提交
// @Param Authorization header string true "Authorization"
// @Param problemIdentity query string true "problemUdentity"
// @Param code body string true "内容"
// @Success 200 {object} SuccessResponse "成功返回问题详情"
// @Failure 400 {object} FailResponse "请求参数错误"
// @Failure 401 {object} FailResponse "未授权"
// @Failure 403 {object} FailResponse "权限不足"
// @Failure 404 {object} FailResponse "资源不存在"
// @Failure 500 {object} FailResponse "服务器内部错误"
// @Router /submissions [post] {}
func Submit(c *gin.Context) {
	problemIdentity := c.DefaultQuery("problemIdentity", "")
	code, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "代码不能为空",
		})
		return
	}

	// 代码保存
	path, err := util.CodeSave(code)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": "代码保存失败: " + err.Error(),
		})
		return
	}

	// 代码提交
	u, _ := c.Get("user")
	userClaim := u.(*util.UserClaims)
	sb := &models.SubmitBasic{
		Identity:        util.GetUUID(),
		ProblemIdentity: problemIdentity,
		UserIdentity:    userClaim.Identity,
		Path:            path,
	}
	// 代码运行判断
	pb := new(models.ProblemBasic)
	err = models.DB.Where("identity=?", problemIdentity).Preload("TestCaseBasics").First(&pb).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "未找到相应题目",
		})
		return
	}
	//  判题状态 channel阻塞
	WA := make(chan int)
	OOM := make(chan int)
	CE := make(chan int, 1)
	passCount := 0
	var maxRunTimeMs int64
	var lock sync.Mutex

	// 错误提示信息
	msg := ""

	// 每次提交只编译一次，编译过程发生在运行时计时开始之前。
	binaryPath := filepath.Join(filepath.Dir(path), "main")
	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}
	var compileStderr bytes.Buffer
	compileCmd := exec.Command("go", "build", "-o", binaryPath, path)
	compileCmd.Stderr = &compileStderr
	compileErr := compileCmd.Run()
	if compileErr != nil {
		msg = compileStderr.String()
		if msg == "" {
			msg = compileErr.Error()
		}
		CE <- 1
	}

	// 运行全部测试样例
	if compileErr == nil {
		for _, testCase := range pb.TestCaseBasics {
			go func(testCase *models.TestCaseBasic) {
				sb.Status = 0
				cmd := exec.Command(binaryPath)
				var out, stderr bytes.Buffer
				cmd.Stderr = &stderr
				cmd.Stdout = &out

				stdinPipe, err := cmd.StdinPipe()
				if err != nil {
					fmt.Println("stdinPipe err")
					msg = err.Error()
					CE <- 1
					return
				}

				runStart := time.Now()
				if err = cmd.Start(); err != nil {
					_ = stdinPipe.Close()
					msg = err.Error()
					CE <- 1
					return
				}

				// 读取运行前内存
				var beginMem runtime.MemStats
				runtime.ReadMemStats(&beginMem)

				if _, err = io.WriteString(stdinPipe, testCase.Input); err != nil {
					fmt.Println("io.WriteString err")
					_ = cmd.Process.Kill()
					_ = cmd.Wait()
					msg = err.Error()
					CE <- 1
					return
				}

				if err = stdinPipe.Close(); err != nil {
					fmt.Println("close stdinPipe err")
				}

				if err = cmd.Wait(); err != nil {
					log.Println("stderr:", stderr.String())
					msg = stderr.String()
					if msg == "" {
						msg = err.Error()
					}
					CE <- 1
					return
				}
				runTimeMs := time.Since(runStart).Milliseconds()
				lock.Lock()
				if runTimeMs > maxRunTimeMs {
					maxRunTimeMs = runTimeMs
				}
				lock.Unlock()

				// 读取运行后内存
				var endMem runtime.MemStats
				runtime.ReadMemStats(&endMem)

				// 答案错误
				if testCase.Output != out.String() {
					msg = " 答案错误"
					WA <- 1
					return
				}

				// 超出内存限制
				if (endMem.Alloc/1024)-(beginMem.Alloc/1024) > uint64(pb.MaxMem) {
					msg = "超出内存限制"
					OOM <- 1
					return
				}
				lock.Lock()
				passCount++
				lock.Unlock()
				// 到这里就算通过
			}(testCase)
		}
	}
	// 阻塞判断
	// [0-待判断 1-答案正确 2-答案错误 3-超出时间限制 4-超出内存限制 5-编译错误]
	select {
	case <-WA:
		sb.Status = 2
	case <-OOM:
		sb.Status = 4
	case <-CE:
		sb.Status = 5
	case <-time.After(time.Millisecond * time.Duration(pb.MaxRuntime)):
		lock.Lock()
		allPassed := passCount != 0 && passCount == len(pb.TestCaseBasics)
		judgingRunTimeMs := maxRunTimeMs
		lock.Unlock()
		if allPassed {
			sb.Status = 1
			msg = fmt.Sprintf("运行时间: %dms", judgingRunTimeMs)
		} else {
			msg = "超出时间限制"
			sb.Status = 3
		}
	}

	// 创建提交记录
	err = models.DB.Create(sb).Error
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": "submit error: " + err.Error(),
		})
		return
	}

	if err = models.DB.Transaction(func(tx *gorm.DB) error {
		var addPassNum int
		if sb.Status == 1 {
			addPassNum = 1
		} else {
			addPassNum = 0
		}

		err = models.DB.Model(&models.UserBasic{}).
			Where("identity = ?", userClaim.Identity).
			Updates(map[string]interface{}{
				"submit_num": gorm.Expr("submit_num + ?", 1),
				"pass_num":   gorm.Expr("pass_num + ?", addPassNum),
			}).Error
		if err != nil {

			return err
		}
		err = models.DB.Model(&models.ProblemBasic{}).
			Where("identity = ?", problemIdentity).
			Updates(map[string]interface{}{
				"submit_num": gorm.Expr("submit_num + ?", 1),
				"pass_num":   gorm.Expr("pass_num + ?", addPassNum),
			}).Error
		if err != nil {

			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": " 更新用户或问题的提交和通过数量出错" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"status": sb.Status,
			"msg":    msg,
		},
		"message": "success",
	})
}
