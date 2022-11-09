package token

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joexu01/ingress-gateway/public"
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

func (c *DefaultTokenClaims) GenerateToken(secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signedString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return signedString, nil
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

func (params *IssueRequest) BindValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
