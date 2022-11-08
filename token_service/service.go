package token

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/gin-gonic/gin"
	cert "github.com/joexu01/ingress-gateway/certificates"
	"github.com/joexu01/ingress-gateway/lib"
	"github.com/joexu01/ingress-gateway/middleware"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var ServiceHttpSrvHandler *http.Server

func IssuerInitRouter() *gin.Engine {
	router := gin.Default()
	router.Use(
		middleware.TranslationMiddleware(),
		middleware.RecoveryMiddleware(),
	)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.POST("/issue", IssueMicroserviceTokenHandler)
	return router
}

func IssueMicroserviceTokenHandler(c *gin.Context) {
	req := &IssueRequest{}
	//err := req.BindValidParams(c)
	err := c.ShouldBindJSON(req)
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	remoteIP := c.RemoteIP()
	log.Println("Remote IP Addr:", remoteIP)

	if remoteIP != req.SourceServiceIP {
		middleware.ResponseError(c, 2001, errors.New("ip addresses didn't match"))
		return
	}

	microserviceToken, err := IssueMicroserviceToken(req)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	middleware.ResponseSuccess(c, microserviceToken)
}

func HttpsServerRun() {
	gin.SetMode(lib.ConfBase.DebugMode)

	// load certs
	caCertFile, err := ioutil.ReadFile(cert.Path(lib.GetStringConf("token_service.https.certs.ca")))
	if err != nil {
		log.Fatalf("error reading CA certificate: %v", err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCertFile)

	r := IssuerInitRouter()
	ServiceHttpSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("token_service.https.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("token_service.https.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("token_service.https.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("token_service.https.max_header_bytes")),
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
			ClientCAs:  certPool,
			MinVersion: tls.VersionTLS12,
		},
	}
	pwd, _ := os.Getwd()
	log.Printf(" [INFO] Present Working Directory: %s\n", pwd)
	go func() {
		log.Printf(" [INFO] HttpsServerRun:%s\n", lib.GetStringConf("token_service.https.addr"))
		if err := ServiceHttpSrvHandler.ListenAndServeTLS(
			cert.Path(lib.GetStringConf("token_service.https.certs.server_cert")), cert.Path(lib.GetStringConf("token_service.https.certs.server_key"))); err != nil {
			log.Fatalf(" [ERROR] HttpsServerRun:%s err:%v\n", lib.GetStringConf("token_service.https.addr"), err)
		}
	}()
}

func HttpsServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := ServiceHttpSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpsServerStop err:%v\n", err)
	}
	log.Printf(" [INFO] HttpsServerStop stopped\n")
}
