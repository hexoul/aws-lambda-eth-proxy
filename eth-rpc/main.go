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
	paramFuncName = "func"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate RPC request
	req := json.GetRpcRequestFromJson(request.Body)
	if method := request.QueryStringParameters[paramFuncName]; method != "" {
		req.Method = method
	} else if method := request.PathParameters[paramFuncName]; method != "" {
		req.Method = method
	}
	fmt.Printf("RpcRequest: %#v\n", req)

	// Forward RPC request to Ether node
	respBody := rpc.DoRpc(rpc.TestnetUrl, req)

	// Relay a response from the node
	resp := json.GetRpcResponseFromJson(respBody)
	fmt.Printf("RpcResponse: %#v\n", resp)
	retCode := 200
	if resp.Error.Code != 0 {
		retCode = 400
	}
	return events.APIGatewayProxyResponse{Body: respBody, StatusCode: retCode}, nil
}

func main() {
	lambda.Start(Handler)
}
