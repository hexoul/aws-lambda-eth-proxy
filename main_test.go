package main

import (
	"testing"

	"github.com/hexoul/aws-lambda-eth-proxy/json"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	req := json.RpcRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getBalance",
		Id:      1,
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
