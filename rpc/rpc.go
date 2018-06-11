// RPC with ethereum node
package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	ethjson "github.com/hexoul/aws-lambda-eth-proxy/json"
)

type Rpc struct {
	netType string
	client  *http.Client
}

const (
	Mainnet = "MAIN"
	Testnet = "TEST"
)

// For singleton
var instance *Rpc
var once sync.Once

// mode is MAINNET or TESTNET
func GetInstance(_netType string) *Rpc {
	once.Do(func() {
		instance = &Rpc{
			netType: _netType,
		}
		instance.InitClient()
	})
	return instance
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

func (r *Rpc) InitClient() {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	r.client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
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
	resp, err := r.client.Post(url, ContentType, reqBody)
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

func (r *Rpc) Call(to, data string) (string, error) {
	req := ethjson.RpcRequest{
		Jsonrpc: "2.0",
		Id:      1,
		Method:  "eth_call",
	}
	params := map[string]string{
		"to":   to,
		"data": data,
	}
	req.Params = append(req.Params, params)
	req.Params = append(req.Params, "latest")
	return r.DoRpc(req)
}
