package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	TokenRequestTypeGateway      = "gateway"
	TokenRequestTypeMicroservice = "microservice"
)

var (
	addr         = "172.16.63.132:3008"
	addrHostname = "172.16.63.132"
	secretStr    = "askdm*6kajsd%^^&asm"
	vegetableIp  = "172.16.63.131"
	potatoIp     = "172.16.63.132"
	//potatoUrl      = "http://172.16.63.132:3008/"
	verifyTokenUrl = "http://gateway-token.io:8881/verify/"
	clientTLS      *http.Client
	client         *http.Client
)

type Response struct {
	ErrorCode int         `json:"errno"`
	ErrorMsg  string      `json:"errmsg"`
	Data      interface{} `json:"data"`
	TraceId   interface{} `json:"trace_id"`
	Stack     interface{} `json:"stack"`
}

type DefaultTokenClaims struct {
	TokenType       string        `json:"ttp"`
	UserID          string        `json:"uid"`
	TargetServiceIP string        `json:"tip"`
	RequestResource string        `json:"rqr"`
	SourceServiceIP string        `json:"sip"`
	GenerateTime    int64         `json:"grt"`
	Context         []ContextItem `json:"ctx"`
	jwt.RegisteredClaims
}

type IssueRequest struct {
	RequestType string `json:"requestType"`

	SourceService   string `json:"sourceService"`
	SourceServiceIP string `json:"sourceServiceIP"`

	TargetService   string `json:"targetService"`
	TargetServiceIP string `json:"targetServiceIP"`
	RequestResource string `json:"requestResource"`

	UserID string `json:"userID"`

	PreviousToken string `json:"previousToken"`
}

type ContextItem struct {
	TargetServiceIP string `json:"tip"`
	RequestResource string `json:"rqr"`
	SourceServiceIP string `json:"sip"`
	GenerateTime    int64  `json:"grt"`
}

func echo(rw http.ResponseWriter, r *http.Request) {
	internal := r.Header.Get("Internal-Token")
	log.Println("Internal Token:", internal)

	gateway := r.Header.Get("Gateway-Token")
	log.Println("Gateway Token:", gateway)

	if ok := verifyGatewayToken(gateway); !ok {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write([]byte("error: failed to verify gateway token"))
		return
	}

	internalToken, err := jwt.ParseWithClaims(internal, &DefaultTokenClaims{}, func(token *jwt.Token) (interface{}, error) { // 解析token
		return []byte(secretStr), nil
	})
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write([]byte("error: failed to parse token string" + err.Error()))
		return
	}

	userClaims, ok := internalToken.Claims.(*DefaultTokenClaims)
	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write([]byte("failed to get token claims"))
		return
	}
	claimStr := fmt.Sprintf("%+v", userClaims)

	var rspText string

	sourceIP := r.RemoteAddr

	rspText = "<h1>" + "Microservice Addr: " + addr + "</h1>"

	rspText += "<h1>" + "Remote Addr: " + sourceIP + "</h1>"

	for k, v := range r.Header {
		h := "<h2>"
		h += k
		h += ": "
		for _, vv := range v {
			h += vv
			h += ""
		}
		h += "</h2>"

		rspText += h
	}

	rspText += claimStr
	rw.Header().Add("Content-Type", "text/html")
	rw.WriteHeader(http.StatusOK)
	log.Println(rspText)
	_, _ = rw.Write([]byte(rspText))
}

func potato(rw http.ResponseWriter, req *http.Request) {
	internal := req.Header.Get("Internal-Token")
	log.Println("Internal Token:", internal)

	gateway := req.Header.Get("Gateway-Token")
	log.Println("Gateway Token:", gateway)

	if ok := verifyGatewayToken(gateway); !ok {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write([]byte("error: failed to verify gateway token"))
		return
	}

	internalToken, err := jwt.ParseWithClaims(internal, &DefaultTokenClaims{}, func(token *jwt.Token) (interface{}, error) { // 解析token
		return []byte(secretStr), nil
	})
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write([]byte("error: failed to parse token string" + err.Error()))
		return
	}

	userClaims, ok := internalToken.Claims.(*DefaultTokenClaims)
	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write([]byte("failed to get token claims"))
		return
	}

	//claimsStr := fmt.Sprintf("%+v", userClaims)
	respStr := fmt.Sprintf(`{"userID":"%s","vegetableType":"potato","amount":3000}`, userClaims.UserID)
	_, _ = rw.Write([]byte(respStr))
}

func requestForInternalToken(req *IssueRequest) (string, error) {
	r, _ := json.Marshal(req)
	reader := bytes.NewReader(r)
	resp, err := clientTLS.Post("https://gateway-token.io:8882/issue", "application/json", reader)
	if err != nil {
		return "", err
	}

	respData := &Response{}
	readAll, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(readAll, respData)

	log.Println("respData:", *respData)

	token, ok := respData.Data.(string)
	if !ok {
		return "", errors.New("failed to get token")
	}

	return token, nil
}

func verifyGatewayToken(gatewayToken string) bool {
	resp, err := client.Get(verifyTokenUrl + gatewayToken)

	if err != nil {
		return false
	}
	verifyResult, _ := ioutil.ReadAll(resp.Body)
	respStruct := &Response{}
	_ = json.Unmarshal(verifyResult, respStruct)

	log.Printf("Response Struct: %+v\n", string(verifyResult))
	ok := respStruct.Data.(bool)
	return ok
}

func main() {
	cert, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		log.Fatalf("could not open certificate file: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	certificate, err := tls.LoadX509KeyPair("client2.crt", "client2.key")
	if err != nil {
		log.Fatalf("could not load certificate: %v", err)
	}

	clientTLS = &http.Client{
		Timeout: time.Minute * 3,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{certificate},
			},
		},
	}

	client = &http.Client{Timeout: time.Minute * 3}

	log.Println("http clientTLS initialized")

	http.HandleFunc("/echo", echo) // 注册路由以及回调函数
	http.HandleFunc("/potato", potato)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
