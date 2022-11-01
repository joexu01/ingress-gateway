package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/joexu01/ingress-gateway/middleware"
	"github.com/joexu01/ingress-gateway/public"
	"github.com/joexu01/ingress-gateway/service"
	"github.com/pkg/errors"
	"strings"
)

func HTTPStripURIMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}

		serviceDetail := serviceInterface.(*service.Detail)
		if serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL && serviceDetail.HTTPRule.NeedStripUri {
			c.Request.URL.Path = strings.Replace(c.Request.URL.Path, serviceDetail.HTTPRule.Rule, "", 1)
		}
	}
}
