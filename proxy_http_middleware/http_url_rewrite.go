package middleware

import (
	"github.com/gin-gonic/gin"
)

func HTTPURLRewriteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//serviceInterface, ok := c.Get("service")
		//if !ok {
		//	middleware.ResponseError(c, 2001, errors.New("service not found"))
		//	c.Abort()
		//	return
		//}
		//
		//serviceDetail := serviceInterface.(*service.Detail)
		//for _, rule := range strings.Split(serviceDetail.HTTPRule.UrlRewrite, ",") {
		//	items := strings.Split(rule, " ")
		//	if len(items) != 2 {
		//		continue
		//	}
		//	compile, err := regexp.Compile(items[0])
		//	if err != nil {
		//		log.Println("regexp.Compile error:", err)
		//		continue
		//	}
		//	replacedPath := compile.ReplaceAll([]byte(c.Request.URL.Path), []byte(items[1]))
		//	c.Request.URL.Path = string(replacedPath)
		//}
		c.Next()
	}
}
