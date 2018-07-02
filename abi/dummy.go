package abi

import (
	"fmt"
	"math/big"

	"github.com/hexoul/aws-lambda-eth-proxy/crypto"
	"github.com/hexoul/aws-lambda-eth-proxy/json"
	"github.com/hexoul/aws-lambda-eth-proxy/rpc"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// DummySendTransaction invokes abi.SendTransaction with dummy of Crypto struct
func DummySendTransaction(abi abi.ABI, to, name string, inputs []interface{}, gas int) (resp json.RPCResponse, err error) {
	data, err := Pack(abi, name, inputs...)
	if err != nil {
		return
	}

	c := crypto.GetDummy()
	rpc.NetType = rpc.Testnet
	r := rpc.GetInstance()
	respStr, err := r.SendTransaction(c.GetAddress(), to, data, gas)
	if err != nil {
		return
	}

	resp = json.GetRPCResponseFromJSON(respStr)
	return
}

// DummySendTransactionWithSign invokes abi.SendTransactionWithSign with dummy of Crypto struct
func DummySendTransactionWithSign(abi abi.ABI, to, name string, inputs []interface{}, gasLimit, gasPrice uint64) (resp json.RPCResponse, err error) {
	data, err := abi.Pack(name, inputs...)
	if err != nil {
		return
	}

	c := crypto.GetDummy()
	rpc.NetType = rpc.Testnet
	r := rpc.GetInstance()

	// Make TX function to get nonce
	tx := func(nonce uint64) (err error) {
		tx := types.NewTransaction(nonce, common.HexToAddress(to), zero, uint64(gasLimit), big.NewInt(int64(gasPrice)), data)
		if tx, err = c.SignTx(tx); err != nil {
			return
		}

		var rlpTx []byte
		if rlpTx, err = rlp.EncodeToBytes(tx); err != nil {
			return
		}

		var respStr string
		if respStr, err = r.SendRawTransaction(rlpTx); err != nil {
			return
		}

		if resp = json.GetRPCResponseFromJSON(respStr); resp.Error == nil {
			return fmt.Errorf("%s", resp.Error.Message)
		}
		return
	}

	c.ApplyNonce(tx)
	return
}
