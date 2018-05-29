package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	ethjson "github.com/hexoul/eth-rpc-on-aws-lambda/eth-rpc/json"
)

func DoRpc(url string, req interface{}) (ret string) {
	// Validate request type
	var msg string
	switch req.(type) {
	case string:
		msg, _ = req.(string)
		break
	case ethjson.RpcRequest:
		if ret, err := json.Marshal(req); err == nil {
			msg = string(ret)
			break
		}
	default:
		return
	}

	// HTTP request
	reqBody := bytes.NewBufferString(msg)
	resp, err := http.Post(url, ContentType, reqBody)
	if err != nil {
		fmt.Printf("DoRpc: HttpError, %s\n", err)
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("DoRpc: IoError, %s\n", err)
		return
	}
	ret = string(respBody)
	resp.Body.Close()
	return
}
