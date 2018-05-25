package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"./rpc"
	"github.com/hexoul/eth-rpc-on-aws-lambda/eth-rpc/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate RPC request
	req := json.GetRpcRequestFromJson(request.Body)
	fmt.Printf("%#v\n", req)

	// Forward RPC request to Ether node
	respBody := rpc.DoRpc(TestnetUrl, request.Body)

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
