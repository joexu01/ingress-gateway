package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joexu01/ingress-gateway/lib"
	"github.com/joexu01/ingress-gateway/middleware"
	"github.com/joexu01/ingress-gateway/public"
	reverse "github.com/joexu01/ingress-gateway/reverse_proxy"
	"github.com/joexu01/ingress-gateway/service"
	token "github.com/joexu01/ingress-gateway/token_service"
	"github.com/joexu01/ingress-gateway/user"
	"github.com/pkg/errors"
	"log"
	"net/url"
)

// HTTPReverseProxyMiddlewareWithoutVerifications HTTP反向代理中间件
func HTTPReverseProxyMiddlewareWithoutVerifications() gin.HandlerFunc {
	return func(c *gin.Context) {
		idTokenStr := c.GetHeader("User-Identity-Token")
		log.Println("User ID Token:", idTokenStr)

		idToken, err := jwt.ParseWithClaims(idTokenStr, &user.DefaultClaims{}, func(token *jwt.Token) (interface{}, error) { // 解析token
			return []byte(lib.GetStringConf("base.jwt.jwt_secret")), nil
		})
		if err != nil {
			middleware.ResponseError(c, 2400, err)
			c.Abort()
			return
		}

		userClaims, ok := idToken.Claims.(*user.DefaultClaims)
		if !ok {
			middleware.ResponseError(c, 2400, errors.New("failed to get token claims"))
			c.Abort()
			return
		}

		// 获取代理服务的详情

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
		nextAddr, err := lb.Get(c.Request.URL.String())
		if err != nil {
			middleware.ResponseError(c, 2003, errors.New("service not found"))
			c.Abort()
			return
		}

		trans, err := service.TransporterHandler.GetTrans(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 2004, err)
			c.Abort()
			return
		}

		// 转换 external / internal Token

		parse, _ := url.Parse(nextAddr)
		nextAddrIp := parse.Hostname()

		tokenReq := &token.IssueRequest{
			RequestType:     public.TokenRequestTypeGateway,
			SourceService:   "Gateway",
			SourceServiceIP: lib.GetStringConf("proxy.http.internal_ip"),
			TargetService:   serviceDetail.Info.ServiceName,
			TargetServiceIP: nextAddrIp,
			RequestResource: c.Request.URL.Path,
			UserID:          userClaims.UserID,
			PreviousToken:   "",
		}

		gatewayToken, err := token.IssueGatewayToken(tokenReq)
		if err != nil {
			middleware.ResponseError(c, 2004, err)
			c.Abort()
			return
		}

		// 向请求头部添加 token

		c.Request.Header.Add("Internal-Token", gatewayToken)
		c.Request.Header.Del("User-Identity-Token")

		//使用reverse proxy.ServeHTTP(c.Request, c.Response)
		rp := reverse.NewLoadBalanceReverseProxy(c, nextAddr, trans)
		rp.ServeHTTP(c.Writer, c.Request)

		c.Next()
	}
}
