package main

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/joexu01/ingress-gateway/lib"
	"log"
)

type DefaultClaims struct {
	UserID   string `json:"userID"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (c *DefaultClaims) GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	ss, err := token.SignedString([]byte(lib.GetStringConf("base.jwt.jwt_secret")))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func main() {
	lib.InitModule("./conf/dev/", []string{"base"})
	defer lib.Destroy()

	log.Println(lib.GetStringConf("base.jwt.jwt_secret"))

	user1 := &DefaultClaims{
		UserID:   "30001",
		Username: "john doe",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Identity Provider",
			Subject:   "User-Identity-Token",
			Audience:  nil,
			ExpiresAt: nil,
			NotBefore: nil,
			IssuedAt:  nil,
			ID:        "",
		},
	}

	user1Token, _ := user1.GenerateToken()

	user2 := &DefaultClaims{
		UserID:   "30002",
		Username: "alice doe",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Identity Provider",
			Subject:   "User-Identity-Token",
			Audience:  nil,
			ExpiresAt: nil,
			NotBefore: nil,
			IssuedAt:  nil,
			ID:        "",
		},
	}

	user2Token, _ := user2.GenerateToken()

	log.Println("user1", user1Token)
	log.Println("user2", user2Token)
}
