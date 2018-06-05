// RPC with ethereum node
package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"

	ethjson "github.com/hexoul/eth-rpc-on-aws-lambda/eth-rpc/json"
)

type Rpc struct {
	netType string
}

const (
	Mainnet = "MAIN"
	Testnet = "TEST"
)

// mode is MAINNET or TESTNET
func New(_netType string) *Rpc {
	return &Rpc{
		netType: _netType,
	}
}

func (r *Rpc) getUrl() (url string) {
	switch r.netType {
	case Mainnet:
		url = MainnetUrls[rand.Intn(len(MainnetUrls))]
		break
	case Testnet:
		url = TestnetUrls[rand.Intn(len(TestnetUrls))]
		break
	}
	return
}

// TODO: Retry when fail, give penalty to low-latency node
func (r *Rpc) DoRpc(req interface{}) (string, error) {
	// Get url following netType
	url := r.getUrl()

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
		return "", fmt.Errorf("Invalid req type")
	}

	// HTTP request
	reqBody := bytes.NewBufferString(msg)
	resp, err := http.Post(url, ContentType, reqBody)
	if err != nil {
		return "", fmt.Errorf("HttpError, %s\n", err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("IoError, %s\n", err)
	}
	ret := string(respBody)
	resp.Body.Close()
	return ret, nil
}
