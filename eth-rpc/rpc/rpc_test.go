package rpc

import (
	"testing"

	"github.com/hexoul/eth-rpc-on-aws-lambda/eth-rpc/json"
)

func TestRpc(t *testing.T) {
	testMsg := "{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[\"a\",1],\"id\":100}"
	targetUrl := "http://13.124.160.186:8545"

	// Test with string param
	if respBody := DoRpc(targetUrl, testMsg); len(respBody) == 0 {
		t.Errorf("Failed to RPC with string")
	}

	// Test with RpcRequest param
	testRpcRequest := json.GetRpcRequestFromJson(testMsg)
	if respBody := DoRpc(targetUrl, testRpcRequest); len(respBody) == 0 {
		t.Errorf("Failed to RPC with RpcRequest")
	}
}
