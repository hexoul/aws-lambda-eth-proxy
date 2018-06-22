package abi

import (
	"math/big"
	"sync/atomic"

	"github.com/hexoul/aws-lambda-eth-proxy/crypto"
	"github.com/hexoul/aws-lambda-eth-proxy/json"
	"github.com/hexoul/aws-lambda-eth-proxy/rpc"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// DummySendTransaction invokes abi.SendTransaction with dummy of Crypto struct
func DummySendTransaction(abi abi.ABI, targetNet, to, name string, inputs []interface{}, gas int) (resp json.RPCResponse, err error) {
	data, err := Pack(abi, name, inputs...)
	if err != nil {
		return
	}

	c := crypto.GetDummy()
	r := rpc.GetInstance(targetNet)
	respStr, err := r.SendTransaction(c.Address, to, data, gas)
	if err != nil {
		return
	}

	resp = json.GetRPCResponseFromJSON(respStr)
	return
}

// DummySendTransactionWithSign invokes abi.SendTransactionWithSign with dummy of Crypto struct
func DummySendTransactionWithSign(abi abi.ABI, targetNet, to, name string, inputs []interface{}, gasLimit, gasPrice uint64) (resp json.RPCResponse, err error) {
	data, err := abi.Pack(name, inputs...)
	if err != nil {
		return
	}

	c := crypto.GetDummy()
	r := rpc.GetInstance(targetNet)
	mutex.Lock()
	defer mutex.Unlock()
	if c.Txnonce == 0 {
		c.Txnonce = r.GetTransactionCount(c.Address)
	}
	nonce := atomic.LoadUint64(&c.Txnonce)
	tx := types.NewTransaction(nonce, common.HexToAddress(to), big.NewInt(0), uint64(gasLimit), big.NewInt(int64(gasPrice)), data)
	tx, err = c.SignTx(tx)
	if err != nil {
		return
	}

	rlpTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return
	}

	respStr, err := r.SendRawTransaction(rlpTx)
	if err != nil {
		return
	}

	resp = json.GetRPCResponseFromJSON(respStr)
	if resp.Error == nil {
		atomic.AddUint64(&c.Txnonce, 1)
	}
	return
}
