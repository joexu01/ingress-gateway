package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	addr      = "172.16.63.131:3008"
	secretStr = "asdjhk8hjasd656*asd"
	potatoUrl = "http://microservice2.io:3008/potato"
	client    *http.Client
)

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
	internal := r.Header.Get("Internal-Token")
	log.Println("Internal Token:", internal)

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

	claimStr := fmt.Sprintf("\n\nToken Claims: %+v\n\n", userClaims)

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
	// 取得 tokens
	internal := req.Header.Get("Internal-Token")
	log.Println("Internal Token:", internal)

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

	// 向 Microservice2 发送请求

	getUrl, _ := url.Parse(potatoUrl)
	r := &http.Request{
		Method: "get",
		URL:    getUrl,
		Header: map[string][]string{},
	}
	r.Header.Add("Internal-Token", internal)

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

func main() {
	client = &http.Client{Timeout: time.Minute * 3}

	log.Println("http client initialized")

	http.HandleFunc("/echo", echo) // 注册路由以及回调函数
	http.HandleFunc("/vegetable", vegetable)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
