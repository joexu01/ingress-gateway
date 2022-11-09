package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	environ := os.Environ()

	for _, env := range environ {
		e := env
		split := strings.Split(e, "=")
		if split[0] == "SEC_VER" {
			fmt.Println(split[1])
		}
	}
}
