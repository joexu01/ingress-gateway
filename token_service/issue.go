package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joexu01/ingress-gateway/public"
	"github.com/joexu01/ingress-gateway/secret"
	"log"
	"math/rand"
	"time"
	"unsafe"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxMask = 1<<6 - 1 // All 1-bits, as many as 6
)

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for 10 characters!
	for i, cache, remain := n-1, src.Int63(), 10; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), 10
		}
		b[i] = letterBytes[int(cache&letterIdxMask)%len(letterBytes)]
		i--
		cache >>= 6
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func IssueGatewayToken(request *IssueRequest) (string, error) {
	if request.RequestType != public.TokenRequestTypeGateway {
		return "", errors.New("bad function call: this function only issues a gateway token")
	}

	claims := &DefaultTokenClaims{
		TokenType:       public.TokenRequestTypeGateway,
		UserID:          request.UserID,
		TargetServiceIP: request.TargetService + "|" + request.TargetServiceIP,
		RequestResource: request.RequestResource,
		SourceServiceIP: request.SourceService + "|" + request.SourceServiceIP,
		GenerateTime:    time.Now().UnixNano(),
		RandomStr:       RandStringBytesMaskImprSrc(8),
		Context:         nil,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Gateway",
			Subject:   "Internal Token",
			Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(305 * time.Second)),
			NotBefore: nil,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "",
		},
	}

	remoteSecret := secret.RemoteSecretHandler.RetrieveSecret(request.TargetServiceIP)
	token, err := claims.GenerateToken([]byte(remoteSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func IssueMicroserviceToken(request *IssueRequest) (string, error) {
	if request.RequestType != public.TokenRequestTypeMicroservice {
		return "", errors.New("bad function call: this function only issues a gateway token")
	}

	prevTokenClaims, err := getTokenClaimsFromTokenStr(request.PreviousToken, secret.RemoteSecretHandler.RetrieveSecret(request.SourceServiceIP))
	if err != nil {
		log.Println("SourceIP:", request.SourceServiceIP, "  Secret:", secret.RemoteSecretHandler.RetrieveSecret(request.SourceServiceIP))
		return "", errors.New("invalid previous token claims")
	}

	ctx := prevTokenClaims.Context
	ctx = append(ctx, ContextItem{
		TargetServiceIP: prevTokenClaims.TargetServiceIP,
		RequestResource: prevTokenClaims.RequestResource,
		SourceServiceIP: prevTokenClaims.SourceServiceIP,
		GenerateTime:    prevTokenClaims.GenerateTime,
	})

	claims := DefaultTokenClaims{
		TokenType:       public.TokenRequestTypeMicroservice,
		UserID:          request.UserID,
		TargetServiceIP: request.TargetService + "|" + request.TargetServiceIP,
		RequestResource: request.RequestResource,
		SourceServiceIP: request.SourceService + "|" + request.SourceServiceIP,
		GenerateTime:    time.Now().UnixNano(),
		RandomStr:       RandStringBytesMaskImprSrc(8),
		Context:         ctx,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Gateway|Token Service",
			Subject:   "Access Token",
			Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(185 * time.Second)),
			NotBefore: nil,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "",
		},
	}

	remoteSecret := secret.RemoteSecretHandler.RetrieveSecret(request.TargetServiceIP)
	token, err := claims.GenerateToken([]byte(remoteSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func getTokenClaimsFromTokenStr(tokenStr, secret string) (*DefaultTokenClaims, error) {
	prevToken, err := jwt.ParseWithClaims(tokenStr, &DefaultTokenClaims{}, func(token *jwt.Token) (interface{}, error) { //  parse token
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	prevTokenClaims, ok := prevToken.Claims.(*DefaultTokenClaims)
	if !ok {
		return nil, errors.New("invalid previous token claims")
	}
	return prevTokenClaims, nil
}
