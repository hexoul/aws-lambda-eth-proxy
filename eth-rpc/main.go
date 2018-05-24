package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      uint32        `json:"id"`
}

type RpcResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      uint32 `json:"id"`
	Result  string `json:"result"`
}

func DoRpc(targetUrl string, msg string) string {
	reqBody := bytes.NewBufferString(msg)
	resp, err := http.Post(targetUrl, "application/json", reqBody)
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
	// Parsing JSON body
	var req RpcRequest
	json.Unmarshal([]byte(request.Body), &req)
	fmt.Printf("%#v\n", req)

	// Forward RPC request to Ether node
	respBody := DoRpc("http://13.124.160.186:8545", request.Body)

	// Relay a response from the node
	var resp RpcResponse
	json.Unmarshal([]byte(respBody), &resp)
	fmt.Printf("%#v\n", resp)
	return events.APIGatewayProxyResponse{Body: respBody, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
