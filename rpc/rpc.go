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

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Rpc struct {
	netType string
	client  *http.Client
}

const (
	Mainnet = "MAIN"
	Testnet = "TEST"
	// For initial request
	initParamJsonrpc = "2.0"
	initParamId      = 1
	// Threshold to classify nodes
	threshold = 10
	// RPC retry count
	retryCnt = 3
)

var (
	// For singleton
	instance *Rpc
	once     sync.Once
	// IP => http fail count
	httpFailCnt = make(map[string]int)
	// netType => available length of IP list
	availLen = make(map[string]int)
)

// mode is MAINNET or TESTNET
func GetInstance(_netType string) *Rpc {
	once.Do(func() {
		instance = &Rpc{}
		instance.InitClient()
		availLen[Mainnet] = len(MainnetUrls)
		availLen[Testnet] = len(TestnetUrls)
	})
	instance.netType = _netType
	return instance
}

func (r *Rpc) getUrl() (url string) {
	switch r.netType {
	case Mainnet:
		url = MainnetUrls[rand.Intn(availLen[Mainnet])]
		break
	case Testnet:
		url = TestnetUrls[rand.Intn(availLen[Testnet])]
		break
	}
	return
}

// refreshUrlList sorts url list to avoid bad nodes
// which is not responsible for our request in the past
func (r *Rpc) refreshUrlList(url string) {
	httpFailCnt[url]++
	if httpFailCnt[url] <= threshold {
		return
	}

	// Pick url list following netType
	var p *[]string
	switch r.netType {
	case Mainnet:
		p = &MainnetUrls
		break
	case Testnet:
		p = &TestnetUrls
		break
	}

	// Pick item will be deleted
	delIdx := -1
	for i, item := range *p {
		if item == url {
			delIdx = i
			break
		}
	}

	// Ignore if this url is already removed or not found
	if delIdx >= availLen[r.netType] || delIdx < 0 {
		return
	}

	// Swap last item and the item will be deleted
	l := availLen[r.netType]
	(*p)[l-1], (*p)[delIdx] = (*p)[delIdx], (*p)[l-1]

	// Decrease available length of url list
	availLen[r.netType]--
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

// Retry when fail, give penalty to low-latency node
func (r *Rpc) DoRpc(req interface{}) (ret string, err error) {
	// Get url following netType
	url := r.getUrl()

	// Validate request type
	var msg string
	switch req.(type) {
	case string:
		msg, _ = req.(string)
		break
	case ethjson.RpcRequest:
		if marshal, e := json.Marshal(req); e == nil {
			msg = string(marshal)
			break
		}
	default:
		err = fmt.Errorf("Invalid req type")
		return
	}

	// HTTP request
	reqBody := bytes.NewBufferString(msg)
	var resp *http.Response
	var respBody []byte
	for i := 0; i < retryCnt; i++ {
		resp, err = r.client.Post(url, ContentType, reqBody)
		if err != nil {
			r.refreshUrlList(url)
			continue
		}
		respBody, err = ioutil.ReadAll(resp.Body)
		if err == nil {
			break
		}
	}
	if len(respBody) == 0 {
		return
	}

	ret = string(respBody)
	resp.Body.Close()
	return
}

func initRpcRequest(method string) ethjson.RpcRequest {
	return ethjson.RpcRequest{
		Jsonrpc: initParamJsonrpc,
		Id:      initParamId,
		Method:  method,
	}
}

func (r *Rpc) Call(to, data string) (string, error) {
	req := initRpcRequest("eth_call")
	params := map[string]string{
		"to":   to,
		"data": data,
	}
	req.Params = append(req.Params, params)
	req.Params = append(req.Params, "latest")
	return r.DoRpc(req)
}

func (r *Rpc) GetCode(addr string) (string, error) {
	req := initRpcRequest("eth_getCode")
	req.Params = append(req.Params, addr)
	req.Params = append(req.Params, "latest")
	return r.DoRpc(req)
}

func (r *Rpc) SendTransaction(from, to, data string, gas int) (string, error) {
	req := initRpcRequest("eth_sendTransaction")
	params := map[string]string{
		"from": from,
		"to":   to,
		"gas":  fmt.Sprintf("0x%x", gas),
		"data": data,
	}
	req.Params = append(req.Params, params)
	return r.DoRpc(req)
}

func (r *Rpc) SendRawTransaction(raw []byte) (string, error) {
	req := initRpcRequest("eth_sendRawTransaction")
	req.Params = append(req.Params, hexutil.Encode(raw))
	return r.DoRpc(req)
}
