package rpc

import (
	"testing"

	"github.com/hexoul/aws-lambda-eth-proxy/json"
)

func TestRefreshUrlList(t *testing.T) {
	r := GetInstance(Testnet)
	initLen := len(TestnetUrls)
	target := TestnetUrls[0]
	for i := 0; i < 30; i++ {
		r.refreshURLList(target)
	}
	if (initLen - 1) != availLen[Testnet] {
		t.Errorf("refreshUrlList is abnormal")
	}
}

func TestCall(t *testing.T) {
	r := GetInstance(Testnet)
	if _, err := r.Call("0x11", "0x123"); err != nil {
		t.Errorf("Failed to RPC Call")
	}
}

func TestRpc(t *testing.T) {
	testMsg := "{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[\"a\",1],\"id\":100}"

	r := GetInstance(Testnet)
	// Test with string param
	if _, err := r.DoRPC(testMsg); err != nil {
		t.Errorf("Failed to RPC with string: %s", err)
	}

	// Test with RpcRequest param
	testRPCRequest := json.GetRPCRequestFromJSON(testMsg)
	if _, err := r.DoRPC(testRPCRequest); err != nil {
		t.Errorf("Failed to RPC with RpcRequest: %s", err)
	}
}

func BenchmarkRpc(b *testing.B) {
	testMsg := "{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[\"a\",1],\"id\":100}"

	r := GetInstance(Testnet)
	for i := 0; i < b.N; i++ {
		r.DoRPC(testMsg)
	}
}
