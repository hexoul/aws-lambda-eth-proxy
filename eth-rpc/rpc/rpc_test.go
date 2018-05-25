package rpc

import (
	"testing"
)

func TestRpc(t *testing.T) {
	testMsg := "{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[\"a\",1],\"id\":100}"
	targetUrl := "http://13.124.160.186:8545"

	respBody := DoRpc(targetUrl, testMsg)
	if len(respBody) == 0 {
		t.Errorf("Failed to RPC")
	}
}
