package cache

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joexu01/ingress-gateway/lib"
	"github.com/joexu01/ingress-gateway/middleware"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	TokenCacheServiceHandler TokenCacheService
	TokenHttpSrvHandler      *http.Server
)

func init() {
	NewTokenCacheService("memory")
}

func NewTokenCacheService(cacheType string) {
	switch cacheType {
	case "memory":
		TokenCacheServiceHandler = NewMemoryCacheService(300)
	default:
		TokenCacheServiceHandler = NewMemoryCacheService(300)
	}
}

func TokenInitRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/verify/:token", VerifyToken)
	return router
}

func VerifyToken(c *gin.Context) {
	token := c.Param("token")

	verify := TokenCacheServiceHandler.Verify(token)
	middleware.ResponseSuccess(c, verify)
}

func HttpServerRun() {
	gin.SetMode(lib.ConfBase.DebugMode)
	r := TokenInitRouter()
	TokenHttpSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("base.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("token_cache.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("token_cache.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("token_cache.http.max_header_bytes")),
	}
	pwd, _ := os.Getwd()
	log.Printf(" [INFO] Present Working Directory: %s\n", pwd)
	go func() {
		log.Printf(" [INFO] HttpServerRun:%s\n", lib.GetStringConf("token_cache.http.addr"))
		if err := TokenHttpSrvHandler.ListenAndServe(); err != nil {
			log.Fatalf(" [ERROR] HttpServerRun:%s err:%v\n", lib.GetStringConf("token_cache.http.addr"), err)
		}
	}()
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := TokenHttpSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpServerStop err:%v\n", err)
	}
	log.Printf(" [INFO] HttpServerStop stopped\n")
}
