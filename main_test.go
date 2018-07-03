package main

import (
	"flag"
	"os"
	"testing"

	"github.com/hexoul/aws-lambda-eth-proxy/crypto"
	"github.com/hexoul/aws-lambda-eth-proxy/json"

	"github.com/aws/aws-lambda-go/events"
)

func TestHelp(t *testing.T) {
	os.Args = os.Args[:1]
	main()
}

func TestEnvMain(t *testing.T) {
	os.Setenv(crypto.IsAwsLambda, "")
	os.Setenv(crypto.Path, "crypto/test/testkey")
	os.Setenv(crypto.Passphrase, "")
	flag.Parse()
	main()
}

func TestArgMain(t *testing.T) {
	os.Setenv(crypto.IsAwsLambda, "")
	os.Args[1] = "crypto/test/testkey"
	os.Args[2] = ""
	main()
}

func TestLambdaHandler(t *testing.T) {
	req := json.RPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getBalance",
		ID:      1,
	}
	req.Params = append(req.Params, "0xeeaf5f87cb85433a0db0fc31863b21d1c8279f7d")
	req.Params = append(req.Params, "latest")
	req.Params = append(req.Params, "ether")
	resp, err := lambdaHandler(nil, events.APIGatewayProxyRequest{
		Body: req.String(),
	})
	if resp.Body == "" || err != nil {
		t.Errorf("Failed to start main")
	}
}
