package json

import (
	"encoding/json"
)

type RpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      int32         `json:"id"`
}

type RpcError struct {
	Code    int32  `json:code"`
	Message string `json:message"`
}

type RpcResponse struct {
	Jsonrpc string                 `json:"jsonrpc"`
	Id      int32                  `json:"id"`
	Result  map[string]interface{} `json:"result"`
	Error   RpcError               `json:"error"`
}

func GetRpcRequestFromJson(msg string) RpcRequest {
	var data RpcRequest
	json.Unmarshal([]byte(msg), &data)
	return data
}

func GetRpcResponseFromJson(msg string) RpcResponse {
	var data RpcResponse
	json.Unmarshal([]byte(msg), &data)
	return data
}
