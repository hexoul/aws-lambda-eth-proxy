package rpc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

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
