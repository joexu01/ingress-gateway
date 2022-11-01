package middleware

import (
	"github.com/gin-gonic/gin"
)

func HTTPHeaderTransformMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//serviceInterface, ok := c.Get("service")
		//if !ok {
		//	middleware.ResponseError(c, 2001, errors.New("service not found"))
		//	c.Abort()
		//	return
		//}
		//serviceDetail := serviceInterface.(*dao.ServiceDetail)
		//for _, rule := range strings.Split(serviceDetail.HTTPRule.HeaderTransform, ",") {
		//	items := strings.Split(rule, " ")
		//	if len(items) != 3 {
		//		continue
		//	}
		//	if items[0] == "add" || items[0] == "edit" {
		//		c.Request.Header.Set(items[1], items[2])
		//	}
		//	if items[0] == "del" {
		//		c.Request.Header.Del(items[1])
		//	}
		//}
		c.Next()
	}
}
