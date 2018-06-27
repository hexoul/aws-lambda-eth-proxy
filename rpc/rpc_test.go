package rpc

import (
	"context"
	"testing"

	"github.com/hexoul/aws-lambda-eth-proxy/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestEthClient(t *testing.T) {
	//client, err := ethclient.Dial("https://rinkeby.infura.io")
	client, err := ethclient.Dial(TestnetUrls[0])
	if err != nil {
		t.Fatalf("%s", err)
	}
	_, err = client.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatalf("%s", err)
	}
}

func BenchmarkEthClient(b *testing.B) {
	client, err := ethclient.Dial(TestnetUrls[0])
	if err != nil {
		b.Fatalf("%s", err)
	}

	address := common.HexToAddress("0x00F912f1F41203DaE29b37fc18db8Dbd3cA9833F")
	errCnt := 0
	for i := 0; i < b.N; i++ {
		_, err := client.BalanceAt(context.Background(), address, nil)
		if err != nil {
			errCnt++
		}
	}
	b.Logf("errCnt %d", errCnt)
}

func BenchmarkHttpClient(b *testing.B) {
	r := GetInstance(Testnet)
	req := json.RPCRequest{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  "eth_getBalance",
	}
	req.Params = append(req.Params, "0x00F912f1F41203DaE29b37fc18db8Dbd3cA9833F")
	req.Params = append(req.Params, "latest")

	errCnt := 0
	for i := 0; i < b.N; i++ {
		_, err := r.DoRPC(req)
		if err != nil {
			errCnt++
		}
	}
	b.Logf("errCnt %d", errCnt)
}

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
