package proxy

import (
	"github.com/gin-gonic/gin"

	proxyMiddleware "github.com/joexu01/ingress-gateway/proxy_http_middleware"
)

func InitRouter(m ...gin.HandlerFunc) *gin.Engine {
	// initialize gin router
	router := gin.Default()
	//router := gin.New()
	// TODO: Use gin.New() instead of gin.Default(). gin.New() doesn't Print log information in console
	router.Use(m...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//OAuth
	//oauth := router.Group("/oauth")
	//oauth.Use(
	//	middleware.TranslationMiddleware(),
	//)

	router.Use(
		proxyMiddleware.HTTPAccessModeMiddleware(),
		proxyMiddleware.HTTPJwtAuthTokenMiddleware(),
		proxyMiddleware.HTTPWhiteListMiddleware(),
		proxyMiddleware.HTTPBlackListMiddleware(),
		proxyMiddleware.HTTPHeaderTransformMiddleware(),
		proxyMiddleware.HTTPStripURIMiddleware(),
		proxyMiddleware.HTTPURLRewriteMiddleware(),
		proxyMiddleware.HTTPReverseProxyMiddleware(),
	)

	return router
}
