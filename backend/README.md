# Gin Gorm 后端

### 设置下载代理
> go env -w GOPROXY=https://goproxy.cn,direct

### 安装Gin
> go get -u github.com/gin-gonic/gin

### 安装Gorm及Postgresql驱动
> go get -u gorm.io/gorm
>
> go get -u gorm.io/driver/postgres
> 
### 获取gin-swagger
> go get -u github.com/swaggo/swag/cmd/swag

### 检查是否成功安装
> swag -v

### (如果没有版本提示)找到外部库，找到github.com/swaggo/swag/cmd/swag/main.go, 右键终端打开
> go install

### (必须)每次更新路由刷新swagger
> swag init

### 安装相关swag库
> go get -u github.com/swaggo/gin-swagger
> go get -u github.com/swaggo/files
> 
## 示例：在对方方法上写上注释
```go
// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200","data":"", "message",""}
// @Router /problem [get]
func GetProblemList(c *gin.Context) {
models.GetProblemList()
c.JSON(http.StatusOK, gin.H{
"code":    http.StatusOK,
"data":    models.GetProblemList(),
"message": "success",
})
}
```

### 路由层配置
#### 顶部引入
```md
_docs "github.com/Anyeling0620/OnlineOJ/backend/docs"
swaggerfiles "github.com/swaggo/files"
ginSwagger "github.com/swaggo/gin-swagger"
```
#### 添加路由
```md
r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
```

### 利用swagger检查端口和测试
> go swagger
#### 访问 ```http://localhost:8080/swagger/index.html```