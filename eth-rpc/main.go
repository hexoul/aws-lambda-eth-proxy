package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

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

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
