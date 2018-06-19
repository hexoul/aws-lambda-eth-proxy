package main

import (
	"context"
	"fmt"

	_ "github.com/hexoul/aws-lambda-eth-proxy/crypto"
	"github.com/hexoul/aws-lambda-eth-proxy/json"
	"github.com/hexoul/aws-lambda-eth-proxy/predefined"
	"github.com/hexoul/aws-lambda-eth-proxy/rpc"
	_ "github.com/hexoul/aws-lambda-eth-proxy/web3"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	ParamFuncName = "func"
	Targetnet     = rpc.Testnet
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate RPC request
	req := json.GetRpcRequestFromJson(request.Body)
	if method := request.QueryStringParameters[ParamFuncName]; method != "" {
		req.Method = method
	} else if method := request.PathParameters[ParamFuncName]; method != "" {
		req.Method = method
	}

	var resp json.RpcResponse
	var err error
	if predefined.Contains(req.Method) {
		// Forward RPC request to predefined function
		resp, err = predefined.Forward(req)
	} else {
		// Forward RPC request to Ether node
		respBody, err := rpc.GetInstance(Targetnet).DoRpc(req)
		if err == nil {
			// Relay a response from the node
			resp = json.GetRpcResponseFromJson(respBody)
		}
	}

	// Preprocessing
	/*
		var unit string
		if req.Method == "eth_getBalance" && len(req.Params) > 2 {
			unit = req.Params[2].(string)
			req.Params = req.Params[:2]
		}
	*/

	// Postprocessing
	/*
		if unit != "" {
			if val, err := web3.FromWei(resp.Result.(string), unit); err == nil {
				resp.Result = val
			}
		}
	*/

	retCode := 200
	if err != nil {
		// In case of server-side RPC fail
		fmt.Println(err.Error())
		resp.Error.Message = err.Error()
		retCode = 400
	} else if resp.Error.Code != 0 {
		// In case of ether-node-side RPC fail
		retCode = 400
	}
	return events.APIGatewayProxyResponse{Body: resp.String(), StatusCode: retCode}, nil
}

func main() {
	lambda.Start(Handler)
}
