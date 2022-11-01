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

func HTTPBlackListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}

		var whiteIpList []string
		var blackIpList []string

		serviceDetail := serviceInterface.(*service.Detail)

		if serviceDetail.AccessControl.WhiteList != "" {
			whiteIpList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}
		if serviceDetail.AccessControl.BlackList != "" {
			blackIpList = strings.Split(serviceDetail.AccessControl.BlackList, ",")
		}
		if serviceDetail.AccessControl.Available && len(whiteIpList) == 0 && len(blackIpList) > 0 {
			if !public.InStringSlice(blackIpList, c.ClientIP()) {
				middleware.ResponseError(c, 3001, errors.New(
					fmt.Sprintf("%s in IP black list", c.ClientIP())))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
