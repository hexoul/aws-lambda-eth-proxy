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

	"github.com/hexoul/aws-lambda-eth-proxy/common"
	"github.com/hexoul/aws-lambda-eth-proxy/crypto"
	ethjson "github.com/hexoul/aws-lambda-eth-proxy/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
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
	// HTTP timeout
	httpTimeout = 5
)

var (
	zero = big.NewInt(0)
	// For singleton
	instance *RPC
	once     sync.Once
	// IP => ethclient
	ethClients = make(map[string]*ethclient.Client)
	// IP => http fail count
	httpFailCnt = make(map[string]int)
	// NetType => available length of IP list
	availLen = make(map[string]int)
	// NetType is either mainnet or testnet
	NetType = Testnet
)

// GetInstance returns the instance of RPC
func GetInstance() *RPC {
	once.Do(func() {
		instance = &RPC{}
		instance.InitClient()
		availLen[Mainnet] = len(MainnetUrls)
		availLen[Testnet] = len(TestnetUrls)

		instance.NetType = NetType
		instance.NetVersion = instance.GetChainID()
		instance.GasPrice = instance.GetGasPrice()

		if c := crypto.GetInstance(); c != nil {
			c.InitChainID(instance.NetVersion)
			c.InitNonce(instance.GetTransactionCount(c.GetAddress()))
		}
	})
	return instance
}

func randomURL(netType string) (url string) {
	switch netType {
	case Mainnet:
		url = MainnetUrls[rand.Intn(availLen[Mainnet])]
		break
	case Testnet:
		url = TestnetUrls[rand.Intn(availLen[Testnet])]
		break
	}
	return
}

func (r *RPC) getURL() string {
	return randomURL(r.NetType)
}

// GetEthClient returns ether client among urls included in target net
func (r *RPC) GetEthClient() *ethclient.Client {
	url := randomURL(r.NetType)
	if ethClients[url] == nil {
		ethClients[url], _ = ethclient.Dial(url)
	}
	return ethClients[url]
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
	netTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Second * httpTimeout,
		}).Dial,
		TLSHandshakeTimeout: time.Second * httpTimeout,
	}
	r.client = &http.Client{
		Timeout:   time.Second * httpTimeout,
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
func (r *RPC) GetChainID() *big.Int {
	req := initRPCRequest("net_version")
	if netVersion, err := r.DoRPC(req); err == nil {
		resp := ethjson.GetRPCResponseFromJSON(netVersion)
		offset, base := common.FindOffsetNBase(resp.Result.(string))
		if chainID, ok := zero.SetString(resp.Result.(string)[offset:], base); ok {
			return chainID
		}
	}
	return nil
}

// GetGasPrice invokes RPC "eth_gasPrice"
func (r *RPC) GetGasPrice() uint64 {
	req := initRPCRequest("eth_gasPrice")
	if gasPrice, err := r.DoRPC(req); err == nil {
		resp := ethjson.GetRPCResponseFromJSON(gasPrice)
		offset, base := common.FindOffsetNBase(resp.Result.(string))
		if uint64GasPrice, ok := zero.SetString(resp.Result.(string)[offset:], base); ok {
			return uint64GasPrice.Uint64()
		}
	}
	return 0
}

// GetTransactionCount invokes RPC "eth_getTransactionCount"
func (r *RPC) GetTransactionCount(addr string) uint64 {
	req := initRPCRequest("eth_getTransactionCount")
	req.Params = append(req.Params, addr)
	req.Params = append(req.Params, "latest")
	if retStr, txCntErr := r.DoRPC(req); txCntErr == nil {
		resp := ethjson.GetRPCResponseFromJSON(retStr)
		offset, base := common.FindOffsetNBase(resp.Result.(string))
		if txNonce, ok := zero.SetString(resp.Result.(string)[offset:], base); ok {
			return txNonce.Uint64()
		}
	}
	return 0
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
