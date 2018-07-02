package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/hexoul/aws-lambda-eth-proxy/crypto"
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
	// IsAwsLambda decides if served as AWS lambda or not
	IsAwsLambda = "AWS_LAMBDA"
)

func handler(req json.RPCRequest) (body string, statusCode int) {
	var resp json.RPCResponse
	var err error
	if predefined.Contains(req.Method) {
		// Forward RPC request to predefined function
		resp, err = predefined.Forward(req)
	} else {
		// Forward RPC request to Ether node
		var respBody string
		if respBody, err = rpc.GetInstance().DoRPC(req); err == nil {
			// Relay a response from the node
			resp = json.GetRPCResponseFromJSON(respBody)
		}
	}

	statusCode = 200
	if err != nil {
		// In case of server-side RPC fail
		fmt.Println(err.Error())
		resp.Error = &json.RPCError{
			Code:    -1,
			Message: err.Error(),
		}
		statusCode = 400
	} else if resp.Error != nil && resp.Error.Code != 0 {
		// In case of ether-node-side RPC fail
		statusCode = 400
	}
	body = resp.String()
	return
}

// lambdaHandler handles APIGatewayProxyRequest as JSON-RPC request
func lambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate RPC request
	req := json.GetRPCRequestFromJSON(request.Body)
	if method := request.QueryStringParameters[ParamFuncName]; method != "" {
		req.Method = method
	} else if method := request.PathParameters[ParamFuncName]; method != "" {
		req.Method = method
	}

	respBody, statusCode := handler(req)
	return events.APIGatewayProxyResponse{Body: respBody, StatusCode: statusCode}, nil
}

// httpHandler handles http.Request as JSON-RPC request
func httpHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	req := json.GetRPCRequestFromJSON(string(b))
	respBody, statusCode := handler(req)
	w.WriteHeader(statusCode)
	w.Write([]byte(respBody))
}

func main() {
	rpc.NetType = Targetnet
	crypto.GetInstance()

	if os.Getenv(IsAwsLambda) != "" {
		fmt.Println("Ready to start Lambda")
		lambda.Start(lambdaHandler)
	} else {
		fmt.Println("Ready to start HTTP/HTTPS")
		http.HandleFunc("/", httpHandler)
		http.ListenAndServe(":8545", nil)
		// http.ListenAndServeTLS()
	}
}
