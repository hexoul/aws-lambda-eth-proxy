package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hexoul/eth-rpc-on-aws-lambda/eth-rpc/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Id      int32         `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type RpcError struct {
	Code    int    `json:code"`
	Message string `json:message"`
}

type RpcResponse struct {
	Jsonrpc string                 `json:"jsonrpc"`
	Id      int32                  `json:"id"`
	Result  map[string]interface{} `json:"result"`
	Error   RpcError               `json:"error"`
}

func DoRpc(targetUrl string, msg string) string {
	reqBody := bytes.NewBufferString(msg)
	resp, err := http.Post(targetUrl, ContentType, reqBody)
	if err != nil {
		fmt.Printf("DoRpc: HttpError, %s\n", err)
		return ""
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("DoRpc: IoError, %s\n", err)
		return ""
	}
	ret := string(respBody)
	resp.Body.Close()
	return ret
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate RPC request
	var req RpcRequest
	json.Unmarshal([]byte(request.Body), &req)
	fmt.Printf("%#v\n", req)

	// Forward RPC request to Ether node
	respBody := DoRpc(TestnetUrl, request.Body)

	// Relay a response from the node
	var resp RpcResponse
	json.Unmarshal([]byte(respBody), &resp)
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
