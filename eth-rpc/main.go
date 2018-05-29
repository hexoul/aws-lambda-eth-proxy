package main

import (
	"context"
	"fmt"

	"github.com/hexoul/eth-rpc-on-aws-lambda/eth-rpc/json"
	"github.com/hexoul/eth-rpc-on-aws-lambda/eth-rpc/rpc"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	PARAM_FUNC_NAME = "func"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate RPC request
	req := json.GetRpcRequestFromJson(request.Body)
	if request.QueryStringParameters[PARAM_FUNC_NAME] != "" {
		req.Method = request.QueryStringParameters[PARAM_FUNC_NAME]
	} else if request.PathParameters[PARAM_FUNC_NAME] != "" {
		req.Method = request.PathParameters[PARAM_FUNC_NAME]
	}
	fmt.Printf("%#v\n", req)

	// Forward RPC request to Ether node
	respBody := rpc.DoRpc(rpc.TestnetUrl, request.Body)

	// Relay a response from the node
	resp := json.GetRpcResponseFromJson(respBody)
	fmt.Printf("%#v\n", resp)
	retCode := 200
	if resp.Error.Code != 0 {
		retCode = 400
	}
	return events.APIGatewayProxyResponse{Body: respBody, StatusCode: retCode}, nil
}

func main() {
	lambda.Start(Handler)
}
