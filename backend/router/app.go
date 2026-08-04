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
	r.GET("/problems/lists", service.GetProblemList)
	r.GET("/problems/details", service.GetProblemDetail)

	// 用户
	r.GET("/users/details", service.GetUserDetail)
	r.POST("/users/login", service.Login)
	r.POST("/users/register", service.Register)
	r.POST("/users/code", service.SendCode)
	// 排行榜
	r.GET("/users/rank", service.GetRankList)

	// 提交记录
	r.GET("/submit/lists", service.GetSubmitList)

	//管理员私有
	r.POST("/problems", middlewares.AuthAdminCheck(), service.ProblemCreate)
	return r
}
