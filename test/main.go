package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      uint32        `json:"id"`
}

func TestJson(msg string) bool {
	var jsonData map[string]interface{}
	json.Unmarshal([]byte(msg), &jsonData)
	fmt.Println(jsonData)

	var data RpcRequest
	json.Unmarshal([]byte(msg), &data)
	fmt.Printf("%#v\n", data)
	return true
}

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

func main() {
	testMsg := "{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[\"a\",1],\"id\":100}"
	fmt.Println(testMsg)

	TestJson(testMsg)
	TestRpc("http://13.124.160.186:8545", testMsg)
}
