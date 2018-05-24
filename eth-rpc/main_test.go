package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hexoul/eth-rpc-on-aws-lambda/eth-rpc/json"
	"github.com/hexoul/eth-rpc-on-aws-lambda/eth-rpc/web3"
)

func TestRpc(targetUrl string, msg string) bool {
	reqBody := bytes.NewBufferString(msg)
	resp, err := http.Post(targetUrl, "application/json", reqBody)
	if err != nil {
		return false
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	fmt.Println(string(respBody))
	resp.Body.Close()
	return true
}

func Test() {
	TestJson()

	TestFromWei()

	TestToWei()

	testMsg := "{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[\"a\",1],\"id\":100}"
	TestRpc("http://13.124.160.186:8545", testMsg)
}
