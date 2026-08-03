package router

import (
	_ "github.com/Anyeling0620/OnlineOJ/backend/docs"
	"github.com/Anyeling0620/OnlineOJ/backend/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/ping", service.Ping)
	r.GET("/problems", service.GetProblemList)

	return r
}
