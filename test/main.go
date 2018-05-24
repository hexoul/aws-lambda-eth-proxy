package main

import (
	"encoding/json"
	"fmt"
)

type RpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      uint32        `json:"id"`
}

func main() {

	testMsg := "{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[\"a\",1],\"id\":100}"
	fmt.Println(testMsg)

	var jsonData map[string]interface{}
	json.Unmarshal([]byte(testMsg), &jsonData)
	fmt.Println(jsonData)

	var data RpcRequest
	json.Unmarshal([]byte(testMsg), &data)
	fmt.Printf("%#v", data)
}
