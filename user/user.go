package user

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/joexu01/ingress-gateway/lib"
)

var DetailServiceHandler *DetailService

func init() {
	DetailServiceHandler = &DetailService{userInfo: make(map[string]*user)}
}

type DetailService struct {
	userInfo map[string]*user
}

type user struct {
	UserID        string
	Username      string
	ExternalToken string
}

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

func (d *DetailService) LoadUserInfo() {
	d.userInfo["30001"] = &user{
		UserID:        "30001",
		Username:      "john doe",
		ExternalToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIzMDAwMSIsInVzZXJuYW1lIjoiam9obiBkb2UiLCJpc3MiOiJJZGVudGl0eSBQcm92aWRlciIsInN1YiI6IlVzZXItSWRlbnRpdHktVG9rZW4ifQ.kM6gXW8W1U1jta5UrGVkKg_MBiJR_IJ-9EgYvuxvZVY",
	}

	d.userInfo["30002"] = &user{
		UserID:        "30002",
		Username:      "alice doe",
		ExternalToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIzMDAwMiIsInVzZXJuYW1lIjoiYWxpY2UgZG9lIiwiaXNzIjoiSWRlbnRpdHkgUHJvdmlkZXIiLCJzdWIiOiJVc2VyLUlkZW50aXR5LVRva2VuIn0.3hNvRdVnzxo_wXwkqevOjJWdL--wBblm6yWX4RSVR1o",
	}
}
