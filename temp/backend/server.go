package main

import (
	"log"
	"net/http"
)

var addr = "127.0.0.1:3008"

func echo(rw http.ResponseWriter, r *http.Request) {
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

func main() {
	http.HandleFunc("/", echo) // 注册路由以及回调函数
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
