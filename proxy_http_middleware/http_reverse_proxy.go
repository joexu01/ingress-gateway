package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/joexu01/ingress-gateway/middleware"
	reverse "github.com/joexu01/ingress-gateway/reverse_proxy"
	"github.com/joexu01/ingress-gateway/service"
	"github.com/pkg/errors"
)

//HTTPReverseProxyMiddleware HTTP反向代理中间件
func HTTPReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceInterface, exists := c.Get("service")
		if !exists {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*service.Detail)
		//创建reverse proxy
		lb, err := service.LoadBalanceHandler.GetLoadBalancer(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 2002, errors.New("service not found"))
			c.Abort()
			return
		}

		trans, err := service.TransporterHandler.GetTrans(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			c.Abort()
			return
		}
		//使用reverse proxy.ServeHTTP(c.Request, c.Response)
		rp := reverse.NewLoadBalanceReverseProxy(c, lb, trans)
		rp.ServeHTTP(c.Writer, c.Request)
		c.Next()
	}
}
