package main

import (
	"context"

	_ "github.com/hexoul/aws-lambda-eth-proxy/crypto"
	"github.com/hexoul/aws-lambda-eth-proxy/json"
	"github.com/hexoul/aws-lambda-eth-proxy/rpc"
	"github.com/hexoul/aws-lambda-eth-proxy/web3"

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

	// Preprocessing
	var unit string
	if req.Method == "eth_getBalance" && len(req.Params) > 2 {
		unit = req.Params[2].(string)
		req.Params = req.Params[:2]
	}

	// Forward RPC request to Ether node
	respBody, err := rpc.GetInstance(Targetnet).DoRpc(req)

	// Relay a response from the node
	resp := json.GetRpcResponseFromJson(respBody)

	// Postprocessing
	if unit != "" {
		if val, err := web3.FromWei(resp.Result.(string), unit); err == nil {
			resp.Result = val
		}
	}

	retCode := 200
	if err != nil {
		// In case of server-side RPC fail
		resp.Error.Message = err.Error()
		respBody = resp.String()
		retCode = 400
	} else if resp.Error.Code != 0 {
		// In case of ether-node-side RPC fail
		retCode = 400
	}
	return events.APIGatewayProxyResponse{Body: respBody, StatusCode: retCode}, nil
}

func main() {
	lambda.Start(Handler)
}
