package rpc

import (
	"testing"

	"github.com/hexoul/eth-rpc-on-aws-lambda/eth-rpc/json"
)

func TestRpc(t *testing.T) {
	testMsg := "{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[\"a\",1],\"id\":100}"

	r := New(Testnet)
	// Test with string param
	if respBody := r.DoRpc(testMsg); len(respBody) == 0 {
		t.Errorf("Failed to RPC with string")
	}

	// Test with RpcRequest param
	testRpcRequest := json.GetRpcRequestFromJson(testMsg)
	if respBody := r.DoRpc(testRpcRequest); len(respBody) == 0 {
		t.Errorf("Failed to RPC with RpcRequest")
	}
}
