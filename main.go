package main

import (
	"context"
	"fmt"

	_ "github.com/hexoul/aws-lambda-eth-proxy/crypto"
	"github.com/hexoul/aws-lambda-eth-proxy/json"
	"github.com/hexoul/aws-lambda-eth-proxy/predefined"
	"github.com/hexoul/aws-lambda-eth-proxy/rpc"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	// ParamFuncName is a name indicating function's
	ParamFuncName = "func"
	// Targetnet indicates target network
	Targetnet = rpc.Testnet
)

// Handler handles APIGatewayProxyRequest as JSON-RPC request
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate RPC request
	req := json.GetRPCRequestFromJSON(request.Body)
	if method := request.QueryStringParameters[ParamFuncName]; method != "" {
		req.Method = method
	} else if method := request.PathParameters[ParamFuncName]; method != "" {
		req.Method = method
	}

	var resp json.RPCResponse
	var err error
	if predefined.Contains(req.Method) {
		// Forward RPC request to predefined function
		resp, err = predefined.Forward(req)
	} else {
		// Forward RPC request to Ether node
		respBody, err := rpc.GetInstance(Targetnet).DoRPC(req)
		if err == nil {
			// Relay a response from the node
			resp = json.GetRPCResponseFromJSON(respBody)
		}
	}

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
	predefined.Targetnet = Targetnet
	lambda.Start(Handler)
}
