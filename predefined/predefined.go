package predefined

import (
	"fmt"

	"github.com/hexoul/aws-lambda-eth-proxy/json"
	_ "github.com/hexoul/aws-lambda-eth-proxy/rpc"
)

func foo(req json.RpcRequest) (json.RpcResponse, error) {
	fmt.Println("foo")
	return json.RpcResponse{}, nil
}

func Forward(path string, req json.RpcRequest) (json.RpcResponse, error) {
	for k, v := range predefinedPaths {
		if k == path {
			return v.(func(json.RpcRequest) (json.RpcResponse, error))(req)
		}
	}
	return json.RpcResponse{}, fmt.Errorf("predefined NOT FOUND")
}

func Contains(path string) bool {
	for k, _ := range predefinedPaths {
		if k == path {
			return true
		}
	}
	return false
}

var predefinedPaths = map[string]interface{}{
	"foo": foo,
}
