package proxy

import (
	"context"
	"github.com/gin-gonic/gin"
	cert "github.com/joexu01/ingress-gateway/certificates"
	"github.com/joexu01/ingress-gateway/lib"
	"github.com/joexu01/ingress-gateway/middleware"
	"log"
	"net/http"
	"time"
)

var (
	HttpProxyHandler  *http.Server
	HttpsProxyHandler *http.Server
)

func HttpProxyRun() {
	gin.SetMode(lib.ConfBase.DebugMode)
	r := InitRouter(
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
	)
	HttpProxyHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.http.max_header_bytes")),
	}

	log.Printf(" [INFO] HttpProxyRun:%s\n", lib.GetStringConf("proxy.http.addr"))
	if err := HttpProxyHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf(" [ERROR] HttpProxyRun:%s err:%v\n", lib.GetStringConf("proxy.http.addr"), err)
	}
}

func HttpProxyStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpProxyHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpProxyStop err:%v\n", err)
	}
	log.Printf(" [INFO] HttpProxyStop %v stopped\n", lib.GetStringConf("proxy.http.addr"))
}

func HttpsProxyRun() {
	gin.SetMode(lib.ConfBase.DebugMode)
	r := InitRouter(
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
	)
	HttpsProxyHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.https.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.https.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.https.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.https.max_header_bytes")),
	}
	log.Printf(" [INFO] HttpsProxyRun:%s\n", lib.GetStringConf("proxy.https.addr"))
	if err := HttpsProxyHandler.ListenAndServeTLS(
		cert.Path("server.crt"), cert.Path("server.key")); err != nil && err != http.ErrServerClosed {
		log.Fatalf(" [ERROR] HttpsProxyRun:%s err:%v\n", lib.GetStringConf("proxy.https.addr"), err)
	}
}

func HttpsProxyStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpProxyHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpsProxyStop err:%v\n", err)
	}
	log.Printf(" [INFO] HttpsProxyStop %v stopped\n", lib.GetStringConf("proxy.https.addr"))
}
