package middlewares

import (
	"fmt"
	"github.com/Anyeling0620/OnlineOJ/backend/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// AuthAdminCheck 验证管理员身份
func AuthAdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			auth = c.GetHeader("token")
		}
		if strings.HasPrefix(auth, "Bearer ") {
			auth = strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		} else {
			auth = strings.TrimSpace(auth)
		}
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "authorization token is missing",
			})
			c.Abort()
			return
		}

		userClaim, err := util.AnalyseToken(auth)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "analyse token err",
			})
			c.Abort()
			return
		}

		if userClaim.IsAdmin != 1 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "unauthorized",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
