package router

import (
	_ "github.com/Anyeling0620/OnlineOJ/backend/docs"
	"github.com/Anyeling0620/OnlineOJ/backend/middlewares"
	"github.com/Anyeling0620/OnlineOJ/backend/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// 公有方法
	r.GET("/ping", service.Ping)

	// 问题
	r.GET("/problems", service.GetProblemList)
	r.GET("/problems/:problemIdentity", service.GetProblemDetail)

	// 用户
	r.GET("/users/:userIdentity", service.GetUserDetail)
	r.POST("/users/login", service.Login)
	r.POST("/users/register", service.Register)
	r.POST("/users/code", service.SendCode)
	// 排行榜
	r.GET("/users/rank", service.GetRankList)

	// 提交记录
	r.GET("/submissions", service.GetSubmitList)

	//管理员私有

	r.GET("/categories", middlewares.AuthAdminCheck(), service.GetCategoryList)
	r.POST("/categories", middlewares.AuthAdminCheck(), service.CategoryCreate)
	r.PUT("/categories/:categoryIdentity", middlewares.AuthAdminCheck(), service.CategoryModify)
	r.DELETE("/categories/:categoryIdentity", middlewares.AuthAdminCheck(), service.CategoryDelete)

	r.POST("/problems", middlewares.AuthAdminCheck(), service.ProblemCreate)
	r.PUT("/problems/:problemIdentity", middlewares.AuthAdminCheck(), service.ProblemModify)

	// 用户私有方法

	r.POST("/submissions", middlewares.AuthUserCheck(), service.Submit)

	return r
}
