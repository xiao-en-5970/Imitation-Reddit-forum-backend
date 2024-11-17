package middlewares

import (
	"bluebell/controller"
	"bluebell/pkg/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func JWTAuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		//zap.L().Info(authHeader)
		if authHeader == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)

			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		fmt.Println(parts[1])
		mc, err := jwt.ParseToken(parts[1])
		fmt.Println(mc)
		if err != nil {

			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		c.Set(controller.CtxUserIDKey, mc.UserID)
		c.Next()
	}
}
