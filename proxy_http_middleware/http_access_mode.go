package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/joexu01/ingress-gateway/middleware"
	"github.com/joexu01/ingress-gateway/service"
)

//HTTPAccessModeMiddleware 基于请求信息匹配接入方式
func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessMode, err := service.ManagerHandler.HTTPAccessMode(c)
		if err != nil {
			middleware.ResponseError(c, 1001, err)
			c.Abort()
			return
		}
		//log.Printf("matched accessMode: %+v", accessMode)
		//log.Println(public.Obj2Json(accessMode))
		c.Set("service", accessMode)
		c.Next()
	}
}
