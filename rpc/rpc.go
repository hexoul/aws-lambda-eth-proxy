// Package rpc invokes JSON-RPC with ethereum node
package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/hexoul/aws-lambda-eth-proxy/crypto"
	ethjson "github.com/hexoul/aws-lambda-eth-proxy/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// RPC is a JSON-RPC manager through HTTP
type RPC struct {
	NetType    string
	NetVersion *big.Int
	client     *http.Client
	GasPrice   uint64
}

const (
	// Mainnet is a const string indicates mainnet
	Mainnet = "MAIN"
	// Testnet is a const string indicates testnet
	Testnet = "TEST"
	// For initial request
	initParamJsonrpc = "2.0"
	initParamID      = 1
	// Threshold to classify nodes
	threshold = 10
	// RPC retry count
	retryCnt = 3
)

var (
	// For singleton
	instance *RPC
	once     sync.Once
	// IP => http fail count
	httpFailCnt = make(map[string]int)
	// NetType => available length of IP list
	availLen = make(map[string]int)
)

// GetInstance returns the instance of Rpc
// _netType should be Mainnet or Testnet
func GetInstance(_netType string) *RPC {
	once.Do(func() {
		instance = &RPC{}
		instance.InitClient()
		availLen[Mainnet] = len(MainnetUrls)
		availLen[Testnet] = len(TestnetUrls)

		bigInt := new(big.Int)
		instance.NetType = _netType
		if netVersion, err := instance.GetChainID(); err == nil {
			resp := ethjson.GetRPCResponseFromJSON(netVersion)
			instance.NetVersion, _ = bigInt.SetString(resp.Result.(string), 10)
			crypto.GetInstance().ChainID = instance.NetVersion
		}

		if gasPrice, err := instance.GetGasPrice(); err == nil {
			resp := ethjson.GetRPCResponseFromJSON(gasPrice)
			uint64GasPrice, _ := bigInt.SetString(resp.Result.(string)[2:], 16)
			instance.GasPrice = uint64GasPrice.Uint64()
		}
	})
	return instance
}

func (r *RPC) getURL() (url string) {
	switch r.NetType {
	case Mainnet:
		url = MainnetUrls[rand.Intn(availLen[Mainnet])]
		break
	case Testnet:
		url = TestnetUrls[rand.Intn(availLen[Testnet])]
		break
	}
	return
}

// refreshURLList sorts url list to avoid bad nodes
// which is not responsible for our request in the past
func (r *RPC) refreshURLList(url string) {
	httpFailCnt[url]++
	if httpFailCnt[url] <= threshold {
		return
	}

	// Pick url list following NetType
	var p *[]string
	switch r.NetType {
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
	if delIdx >= availLen[r.NetType] || delIdx < 0 {
		return
	}

	// Swap last item and the item will be deleted
	l := availLen[r.NetType]
	(*p)[l-1], (*p)[delIdx] = (*p)[delIdx], (*p)[l-1]

	// Decrease available length of url list
	availLen[r.NetType]--
}

// InitClient initializes HTTP client to reduce handshaking overhead
func (r *RPC) InitClient() {
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

// DoRPC invokes HTTP post request to ethereum node
// Retry when fail, give penalty to low-latency node
func (r *RPC) DoRPC(req interface{}) (ret string, err error) {
	// Get url following NetType
	url := r.getURL()

	// Validate request type
	var msg string
	switch req.(type) {
	case string:
		msg, _ = req.(string)
		break
	case ethjson.RPCRequest:
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
			r.refreshURLList(url)
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

func initRPCRequest(method string) ethjson.RPCRequest {
	return ethjson.RPCRequest{
		Jsonrpc: initParamJsonrpc,
		ID:      initParamID,
		Method:  method,
	}
}

// Call invokes RPC "eth_call"
func (r *RPC) Call(to, data string) (string, error) {
	req := initRPCRequest("eth_call")
	params := map[string]string{
		"to":   to,
		"data": data,
	}
	req.Params = append(req.Params, params)
	req.Params = append(req.Params, "latest")
	return r.DoRPC(req)
}

// GetCode invokes RPC "eth_getCode"
func (r *RPC) GetCode(addr string) (string, error) {
	req := initRPCRequest("eth_getCode")
	req.Params = append(req.Params, addr)
	req.Params = append(req.Params, "latest")
	return r.DoRPC(req)
}

// GetChainID invokes RPC "net_version"
func (r *RPC) GetChainID() (string, error) {
	req := initRPCRequest("net_version")
	return r.DoRPC(req)
}

// GetGasPrice invokes RPC "eth_gasPrice"
func (r *RPC) GetGasPrice() (string, error) {
	req := initRPCRequest("eth_gasPrice")
	return r.DoRPC(req)
}

// GetTransactionCount invokes RPC "eth_getTransactionCount"
func (r *RPC) GetTransactionCount(addr string) (string, error) {
	req := initRPCRequest("eth_getTransactionCount")
	req.Params = append(req.Params, addr)
	req.Params = append(req.Params, "latest")
	return r.DoRPC(req)
}

// SendTransaction invokes RPC "eth_sendTransaction"
func (r *RPC) SendTransaction(from, to, data string, gas int) (string, error) {
	req := initRPCRequest("eth_sendTransaction")
	params := map[string]string{
		"from": from,
		"to":   to,
		"gas":  fmt.Sprintf("0x%x", gas),
		"data": data,
	}
	req.Params = append(req.Params, params)
	return r.DoRPC(req)
}

// SendRawTransaction invokes RPC "eth_sendRawTransaction"
func (r *RPC) SendRawTransaction(raw []byte) (string, error) {
	req := initRPCRequest("eth_sendRawTransaction")
	req.Params = append(req.Params, hexutil.Encode(raw))
	return r.DoRPC(req)
}
