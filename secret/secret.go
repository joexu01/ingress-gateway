package secret

import (
	"github.com/joexu01/ingress-gateway/lib"
	"time"
)

var RemoteSecretHandler *Service

func init() {
	RemoteSecretHandler = &Service{secretMap: make(map[string]*Item)}
	conf := lib.GetStringMapStringConf("secret.secrets")

	for serviceIP, secretStr := range conf {
		item := &Item{
			SecretStr: secretStr,
			CreatedAt: time.Now().Unix(),
		}
		RemoteSecretHandler.secretMap[serviceIP] = item
	}
}

type Item struct {
	SecretStr string
	CreatedAt int64
}

type Service struct {
	secretMap map[string]*Item
}

func (s *Service) RetrieveSecret(remoteIP string) string {
	if item, ok := s.secretMap[remoteIP]; ok {
		return item.SecretStr
	}
	return ""
}
