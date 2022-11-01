package main

import (
	"github.com/joexu01/ingress-gateway/lib"
	proxy "github.com/joexu01/ingress-gateway/proxy_http_router"
	"github.com/joexu01/ingress-gateway/service"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	_ = lib.InitModule("./conf/dev/", []string{"base", "proxy"})
	defer lib.Destroy()

	_ = service.ManagerHandler.LoadOnce()

	proxy.HttpProxyRun()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	proxy.HttpsProxyStop()
}
