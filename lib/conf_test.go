package lib

import (
	"log"
	"testing"
)

func TestGetStringMapStringConf(t *testing.T) {
	_ = InitModule("../conf/dev/", []string{"base", "secret"})
	defer Destroy()
	conf := GetStringMapStringConf("secret.secrets")
	log.Printf("Secret Map: %+v\n", conf)
}
