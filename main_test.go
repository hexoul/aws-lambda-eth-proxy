package main

import (
	"flag"
	"os"
	"testing"

	"github.com/hexoul/aws-lambda-eth-proxy/crypto"
	"github.com/hexoul/aws-lambda-eth-proxy/json"

	"github.com/aws/aws-lambda-go/events"
)

func TestMain(t *testing.T) {
	os.Setenv(crypto.Passphrase, "")
	os.Setenv(crypto.Path, "crypto/test/testkey")
	flag.Parse()
	main()
}

func TestHandler(t *testing.T) {
	req := json.RPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getBalance",
		ID:      1,
	}
	req.Params = append(req.Params, "0xeeaf5f87cb85433a0db0fc31863b21d1c8279f7d")
	req.Params = append(req.Params, "latest")
	req.Params = append(req.Params, "ether")
	resp, err := Handler(nil, events.APIGatewayProxyRequest{
		Body: req.String(),
	})
	if resp.Body == "" || err != nil {
		t.Errorf("Failed to start main")
	}
}
