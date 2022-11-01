package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joexu01/ingress-gateway/middleware"
	"github.com/joexu01/ingress-gateway/public"
	"github.com/joexu01/ingress-gateway/service"
	"github.com/pkg/errors"
	"strings"
)

func HTTPWhiteListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}

		var ipList []string
		serviceDetail := serviceInterface.(*service.Detail)
		if serviceDetail.AccessControl.WhiteList != "" {
			ipList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}
		if serviceDetail.AccessControl.Available && len(ipList) > 0 {
			if !public.InStringSlice(ipList, c.ClientIP()) {
				middleware.ResponseError(c, 3001, errors.New(
					fmt.Sprintf("%s not in IP white list", c.ClientIP())))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
