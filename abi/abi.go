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

const Targetnet = rpc.Testnet

func Pack(abi abi.ABI, name string, args ...interface{}) (string, error) {
	data, err := abi.Pack(name, args...)
	if err != nil {
		return "", err
	} else {
		return hexutil.Encode(data), nil
	}
}

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

func Call(abi abi.ABI, to, name string, inputs []interface{}, outputs interface{}) (resp json.RpcResponse, err error) {
	data, err := Pack(abi, name, inputs...)
	if err != nil {
		return
	}

	r := rpc.GetInstance(Targetnet)
	respStr, err := r.Call(to, data)
	if err != nil {
		return
	}

	resp = json.GetRpcResponseFromJson(respStr)
	return
}

func SendTransaction(abi abi.ABI, to, name string, inputs []interface{}, gas int) (resp json.RpcResponse, err error) {
	data, err := Pack(abi, name, inputs...)
	if err != nil {
		return
	}

	c := crypto.GetInstance()
	r := rpc.GetInstance(Targetnet)
	respStr, err := r.SendTransaction(c.Address, to, data, gas)
	if err != nil {
		return
	}

	resp = json.GetRpcResponseFromJson(respStr)
	return
}

func SendTransactionWithSign(abi abi.ABI, to, name string, inputs []interface{}, gasLimit, gasPrice uint64) (resp json.RpcResponse, err error) {
	data, err := abi.Pack(name, inputs...)
	if err != nil {
		return
	}

	c := crypto.GetInstance()
	tx := types.NewTransaction(0, common.HexToAddress(to), big.NewInt(0), uint64(gasLimit), big.NewInt(int64(gasPrice)), data)
	tx, err = c.SignTx(tx)
	if err != nil {
		return
	}

	rlpTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return
	}

	r := rpc.GetInstance(Targetnet)
	respStr, err := r.SendRawTransaction(rlpTx)
	if err != nil {
		return
	}

	resp = json.GetRpcResponseFromJson(respStr)
	return
}

func GetAbiFromJson(raw string) (abi.ABI, error) {
	return abi.JSON(strings.NewReader(raw))
}

// getAbiFromAddress is NOT YET SUPPORTED
// TODO: use eth.compile.solidity?
func getAbiFromAddress(addr string) (abi abi.ABI) {
	r := rpc.GetInstance(Targetnet)
	respStr, err := r.GetCode(addr)
	if err != nil {
		return
	}

	json.GetRpcResponseFromJson(respStr)
	return
}
