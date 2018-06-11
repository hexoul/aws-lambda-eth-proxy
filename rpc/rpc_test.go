package rpc

import (
	"testing"

	"github.com/hexoul/eth-rpc-on-aws-lambda/json"
)

func TestRpc(t *testing.T) {
	testMsg := "{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[\"a\",1],\"id\":100}"

	r := GetInstance(Testnet)
	// Test with string param
	if _, err := r.DoRpc(testMsg); err != nil {
		t.Errorf("Failed to RPC with string: %s", err)
	}

	// Test with RpcRequest param
	testRpcRequest := json.GetRpcRequestFromJson(testMsg)
	if _, err := r.DoRpc(testRpcRequest); err != nil {
		t.Errorf("Failed to RPC with RpcRequest: %s", err)
	}
}

func BenchmarkRpc(b *testing.B) {
	testMsg := "{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[\"a\",1],\"id\":100}"

	r := GetInstance(Testnet)
	for i := 0; i < b.N; i++ {
		r.DoRpc(testMsg)
	}
}
