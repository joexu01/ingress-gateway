package token

import (
	"github.com/joexu01/ingress-gateway/lib"
	"github.com/joexu01/ingress-gateway/public"
	"github.com/joexu01/ingress-gateway/secret"
	"testing"
)

func TestIssueGatewayToken(t *testing.T) {
	lib.InitModule("../conf/dev/", []string{"base", "secret", "token_service"})
	defer lib.Destroy()

	secret.RemoteSecretHandler.LoadSecrets()

	req := IssueRequest{
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

}
