// Package abi implements smart contract call helper
package abi

import (
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/hexoul/aws-lambda-eth-proxy/crypto"
	"github.com/hexoul/aws-lambda-eth-proxy/json"
	"github.com/hexoul/aws-lambda-eth-proxy/rpc"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

var zero = big.NewInt(0)

// Pack makes packed data with inputs on ABI
func Pack(abi abi.ABI, name string, args ...interface{}) (string, error) {
	data, err := abi.Pack(name, args...)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(data), nil
}

// Unpack fills output into given ABI
func Unpack(abi abi.ABI, v interface{}, name string, output string) error {
	var data []byte
	var err error
	if output[:2] == "0x" {
		data, err = hex.DecodeString(output[2:])
	} else {
		data, err = hex.DecodeString(output)
	}

	if err != nil {
		return err
	}
	return abi.Unpack(v, name, data)
}

// Call gets contract value with contract address and name
func Call(abi abi.ABI, targetNet, to, name string, inputs []interface{}) (resp json.RPCResponse, err error) {
	data, err := Pack(abi, name, inputs...)
	if err != nil {
		return
	}

	r := rpc.GetInstance(targetNet)
	respStr, err := r.Call(to, data)
	if err != nil {
		return
	}

	resp = json.GetRPCResponseFromJSON(respStr)
	return
}

// SendTransaction calls smart contract with ABI using eth_sendTransaction
func SendTransaction(abi abi.ABI, targetNet, to, name string, inputs []interface{}, gas int) (resp json.RPCResponse, err error) {
	data, err := Pack(abi, name, inputs...)
	if err != nil {
		return
	}

	c := crypto.GetInstance()
	r := rpc.GetInstance(targetNet)
	respStr, err := r.SendTransaction(c.Address, to, data, gas)
	if err != nil {
		return
	}

	resp = json.GetRPCResponseFromJSON(respStr)
	return
}

// SendTransactionWithSign calls smart contract with ABI using eth_sendRawTransaction
func SendTransactionWithSign(abi abi.ABI, targetNet, to, name string, inputs []interface{}, gasLimit, gasPrice uint64) (resp json.RPCResponse, err error) {
	data, err := abi.Pack(name, inputs...)
	if err != nil {
		return
	}

	// Get TX nonce
	c := crypto.GetInstance()
	r := rpc.GetInstance(targetNet)
	if c.Txnonce == 0 {
		c.Txnonce = r.GetTransactionCount(c.Address)
	}

	// TODO: nonce atomic increase
	tx := types.NewTransaction(c.Txnonce, common.HexToAddress(to), zero, uint64(gasLimit), big.NewInt(int64(gasPrice)), data)
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
	return
}

// GetAbiFromJSON returns ABI object from JSON string
func GetAbiFromJSON(raw string) (abi.ABI, error) {
	return abi.JSON(strings.NewReader(raw))
}

// getAbiFromAddress is NOT YET SUPPORTED
// TODO: use eth.compile.solidity?
func getAbiFromAddress(targetNet, addr string) (abi abi.ABI) {
	r := rpc.GetInstance(targetNet)
	respStr, err := r.GetCode(addr)
	if err != nil {
		return
	}

	json.GetRPCResponseFromJSON(respStr)
	return
}
