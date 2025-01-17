package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
)

var (
	addr      = "172.16.63.132:3008"
	secretStr = "askdm*6kajsd%^^&asm"
)

type DefaultTokenClaims struct {
	TokenType       string        `json:"ttp"`
	UserID          string        `json:"uid"`
	TargetServiceIP string        `json:"tip"`
	RequestResource string        `json:"rqr"`
	SourceServiceIP string        `json:"sip"`
	GenerateTime    int64         `json:"grt"`
	RandomStr       string        `json:"rds,omitempty"`
	Context         []ContextItem `json:"ctx"`
	jwt.RegisteredClaims
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

	log.Printf("\n\nToken Claims: %+v\n\n", userClaims)

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

	rw.Header().Add("Content-Type", "text/html")
	rw.WriteHeader(http.StatusOK)
	log.Println(rspText)
	_, _ = rw.Write([]byte(rspText))
}

func potato(rw http.ResponseWriter, req *http.Request) {
	internal := req.Header.Get("Internal-Token")
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

	log.Printf("\n\nToken Claims: %+v\n\n", userClaims)
	respStr := fmt.Sprintf(`{"userID":"%s","vegetableType":"potato","amount":3000}`, userClaims.UserID)
	_, _ = rw.Write([]byte(respStr))
}

func main() {
	http.HandleFunc("/echo", echo) // 注册路由以及回调函数
	http.HandleFunc("/potato", potato)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
