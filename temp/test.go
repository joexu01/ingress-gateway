package main

import (
	"fmt"
	"github.com/joexu01/ingress-gateway/public"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	fi, _ := os.Open(public.ServiceRuleFile)
	bytes, err := ioutil.ReadAll(fi)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(bytes))

}
