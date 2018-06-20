// Predefined functions for RPC request
package predefined

import (
	"fmt"

	"github.com/hexoul/aws-lambda-eth-proxy/json"
	"github.com/hexoul/aws-lambda-eth-proxy/rpc"
	"github.com/hexoul/aws-lambda-eth-proxy/web3"
)

const Targetnet = rpc.Testnet

func foo(req json.RPCRequest) (json.RPCResponse, error) {
	fmt.Println("foo")
	return json.RPCResponse{}, nil
}

func getBalance(req json.RPCRequest) (json.RPCResponse, error) {
	// Preprocessing
	var unit string
	if len(req.Params) > 2 {
		unit = req.Params[2].(string)
		req.Params = req.Params[:2]
	}

	// RPC
	var resp json.RPCResponse
	respBody, err := rpc.GetInstance(Targetnet).DoRpc(req)
	if err == nil {
		resp = json.GetRPCResponseFromJSON(respBody)
		// Postprocessing
		if unit != "" {
			if val, err := web3.FromWei(resp.Result.(string), unit); err == nil {
				resp.Result = val
			}
		}
	}

	return resp, err
}

func Forward(req json.RPCRequest) (json.RPCResponse, error) {
	for k, v := range predefinedPaths {
		if k == req.Method {
			return v.(func(json.RPCRequest) (json.RPCResponse, error))(req)
		}
	}
	return json.RPCResponse{}, fmt.Errorf("predefined NOT FOUND")
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
	"foo":            foo,
	"eth_getBalance": getBalance,
}
