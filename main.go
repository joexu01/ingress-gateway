package main

import (
	"github.com/joexu01/ingress-gateway/lib"
	proxy "github.com/joexu01/ingress-gateway/proxy_http_router"
	"github.com/joexu01/ingress-gateway/service"
	cache "github.com/joexu01/ingress-gateway/token_cache"
	token "github.com/joexu01/ingress-gateway/token_service"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	_ = lib.InitModule("./conf/dev/", []string{"base", "proxy", "secret", "token_cache"})
	defer lib.Destroy()

	_ = service.ManagerHandler.LoadOnce()

	go proxy.HttpProxyRun()
	go cache.HttpServerRun()
	go token.HttpsServerRun()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	proxy.HttpProxyStop()
	cache.HttpServerStop()
	token.HttpsServerStop()
}
