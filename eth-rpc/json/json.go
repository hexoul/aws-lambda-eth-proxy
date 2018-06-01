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
	Jsonrpc string      `json:"jsonrpc"`
	Id      int32       `json:"id"`
	Result  interface{} `json:"result"`
	Error   RpcError    `json:"error"`
}

func GetRpcRequestFromJson(msg string) RpcRequest {
	var data RpcRequest
	json.Unmarshal([]byte(msg), &data)
	return data
}

func (r *RpcRequest) String() string {
	ret, err := json.Marshal(r)
	if err == nil {
		return string(ret)
	}
	return ""
}

func GetRpcResponseFromJson(msg string) RpcResponse {
	var data RpcResponse
	json.Unmarshal([]byte(msg), &data)
	return data
}

func (r *RpcResponse) String() string {
	ret, err := json.Marshal(r)
	if err == nil {
		return string(ret)
	}
	return ""
}
