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
	"net/url"
	"time"
)

const (
	TokenRequestTypeMicroservice = "microservice"
)

var (
	addr           = "172.16.63.131:3008"
	addrHostname   = "172.16.63.131"
	secretStr      = "asdjhk8hjasd656*asd"
	potatoIp       = "172.16.63.132"
	potatoUrl      = "http://microservice2.io:3008/potato"
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

type VegetableItemResp struct {
	UserID string `json:"userID"`
	Type   string `json:"vegetableType"`
	Amount int    `json:"amount"`
}

type VegetableItem struct {
	Type   string `json:"vegetableType"`
	Amount int    `json:"amount"`
}

type VegetableList struct {
	UserID string `json:"userID"`
	List   []VegetableItem
}

func echo(rw http.ResponseWriter, r *http.Request) {
	handlerPattern := "/echo"

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

	// 这一步是检查到来的请求的 url path 是否与 gateway 看到的请求 url path 一致
	//
	// 一般来说，一个 url 的形式如下
	//
	//	[scheme:][//[userinfo@]host][/]path[?query][#fragment]
	//
	// 比如用户请求的是 /echo ，Gateway 会将这个 path 记录到 Gateway-Token 中
	//
	// 微服务在收到请求后，会检查 Gateway-Token 中的 path 是否与 url 中的 path 一致
	// 如果不一致，说明可能是攻击者利用其他 Gateway-Token 发送的伪造请求

	claimStr := fmt.Sprintf("\n\nToken Claims: %+v\n\n", userClaims)
	if handlerPattern != userClaims.RequestResource {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write([]byte("bad request path"))
	}

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

func vegetable(rw http.ResponseWriter, req *http.Request) {
	handlerPattern := "/vegetable"

	// 取得 tokens
	internal := req.Header.Get("Internal-Token")
	log.Println("Internal Token:", internal)

	gateway := req.Header.Get("Gateway-Token")
	log.Println("Gateway Token:", gateway)

	// 验证 Gateway-Token 的有效性
	if ok := verifyGatewayToken(gateway); !ok {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write([]byte("error: failed to verify gateway token"))
		return
	}

	// 解析 token
	internalToken, err := jwt.ParseWithClaims(internal, &DefaultTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
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

	log.Printf("\n\nToken Claims: %+v\n\n", userClaims)
	if handlerPattern != userClaims.RequestResource {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write([]byte("bad request path"))
	}

	// 请求访问 Microservice2 的 Token
	issueReq := &IssueRequest{
		RequestType:     TokenRequestTypeMicroservice,
		SourceService:   "Vegetable",
		SourceServiceIP: addrHostname,
		TargetService:   "Potato",
		TargetServiceIP: potatoIp,
		RequestResource: "/potato",
		UserID:          userClaims.UserID,
		PreviousToken:   internal,
	}

	log.Printf("issue request: %+v\n", issueReq)

	token, err := requestForInternalToken(issueReq)
	if err != nil {
		log.Println("req for token error:", err)
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write([]byte("req for token error:" + err.Error()))
		return
	}

	log.Printf("new token: %s\n", token[:12])

	// 向 Microservice2 发送请求

	getUrl, _ := url.Parse(potatoUrl)
	r := &http.Request{
		Method: "get",
		URL:    getUrl,
		Header: map[string][]string{},
	}
	r.Header.Add("Gateway-Token", gateway)
	r.Header.Add("Internal-Token", token)

	response, err := client.Do(r)
	if err != nil {
		log.Println("req for token error:", err)
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write([]byte("request for microservice2 error:" + err.Error()))
		return
	}

	all, _ := ioutil.ReadAll(response.Body)

	item := &VegetableItemResp{}
	_ = json.Unmarshal(all, item)

	list := &VegetableList{
		UserID: userClaims.UserID,
		List: []VegetableItem{{
			Type:   item.Type,
			Amount: item.Amount,
		}},
	}

	marshal, _ := json.Marshal(list)

	_, _ = rw.Write([]byte("Response: " + string(marshal)))
	return
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

	certificate, err := tls.LoadX509KeyPair("client1.crt", "client1.key")
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
	http.HandleFunc("/vegetable", vegetable)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
