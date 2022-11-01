package middleware

import (
	"github.com/gin-gonic/gin"
)

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//session := sessions.Default(c)
		//if adminInfo, ok := session.
		//	Get(public.AdminSessionInfoKey).(string); !ok || adminInfo == "" {
		//	ResponseError(c, http.StatusUnauthorized, errors.New("user not login"))
		//	c.Abort()
		//	return
		//}
		c.Next()
	}
}
