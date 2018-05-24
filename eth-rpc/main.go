package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      uint32        `json:"id"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Received body: ", request.Body)
	fmt.Println("Received header len: ", len(request.Headers))

	/*
		for key, value := range request.Headers {
		}
	*/

	// Parsing JSON body
	var jsonData map[string]interface{}
	json.Unmarshal([]byte(request.Body), &jsonData)
	fmt.Println(jsonData)

	var data RpcRequest
	json.Unmarshal([]byte(request.Body), &data)
	fmt.Printf("%#v\n", data)

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
