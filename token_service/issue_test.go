package token

import (
	"github.com/joexu01/ingress-gateway/lib"
	"github.com/joexu01/ingress-gateway/public"
	"github.com/joexu01/ingress-gateway/secret"
	"testing"
)

func TestIssueGatewayToken(t *testing.T) {
	_ = lib.InitModule("../conf/dev/", []string{"base", "secret", "token_service"})
	defer lib.Destroy()

	secret.RemoteSecretHandler.LoadSecrets()

	req := &IssueRequest{
		RequestType:     public.TokenRequestTypeGateway,
		SourceService:   "Gateway",
		SourceServiceIP: "172.16.63.1",
		TargetService:   "Vegetable",
		TargetServiceIP: "172.16.63.131",
		RequestResource: "/list",
		UserID:          "300",
		PreviousToken:   "",
	}

	gatewayToken, err := IssueGatewayToken(req)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(gatewayToken)
}

func TestIssueMicroserviceToken(t *testing.T) {
	_ = lib.InitModule("../conf/dev/", []string{"base", "secret", "token_service"})
	defer lib.Destroy()

	secret.RemoteSecretHandler.LoadSecrets()
	req := &IssueRequest{
		RequestType:     public.TokenRequestTypeMicroservice,
		SourceService:   "Vegetables",
		SourceServiceIP: "172.16.63.131",
		TargetService:   "Potatoes",
		TargetServiceIP: "172.16.63.132",
		RequestResource: "/potatoes",
		UserID:          "300",
		PreviousToken:   `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0dHAiOiJnYXRld2F5IiwidWlkIjoiMzAwIiwidGlwIjoiVmVnZXRhYmxlfDE3Mi4xNi42My4xMzEiLCJycXIiOiIvbGlzdCIsInNpcCI6IkdhdGV3YXl8MTcyLjE2LjYzLjEiLCJncnQiOjE2Njc3ODIzOTEsImN0eCI6bnVsbCwiaXNzIjoiR2F0ZXdheSIsInN1YiI6IkludGVybmFsIFRva2VuIiwiZXhwIjoxNjY3NzgyNjk2LCJpYXQiOjE2Njc3ODIzOTF9.eAnG6TTRSdB12gVcFACp7Gq6zcn2Ipuii2e7hCo29nc`,
	}

	microserviceToken, err := IssueMicroserviceToken(req)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Microservice Token: %s\n", microserviceToken)

	claims, err := getTokenClaimsFromTokenStr(microserviceToken, secret.RemoteSecretHandler.RetrieveSecret("172.16.63.132"))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Claims: %+v\n", claims)
}
