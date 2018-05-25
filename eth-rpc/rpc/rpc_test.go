package rpc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func TestRpc(targetUrl string, msg string) {
	testMsg := "{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[\"a\",1],\"id\":100}"
	targetUrl := "http://13.124.160.186:8545"

	reqBody := bytes.NewBufferString(msg)
	resp, err := http.Post(targetUrl, "application/json", reqBody)
	if err != nil {
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println(string(respBody))
	resp.Body.Close()
}
