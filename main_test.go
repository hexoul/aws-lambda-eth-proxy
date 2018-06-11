package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	body := "{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[\"a\",1],\"id\":100}"
	resp, err := Handler(nil, events.APIGatewayProxyRequest{
		Body: body,
	})
	if resp.Body == "" || err != nil {
		t.Errorf("Failed to start main")
	}
}
