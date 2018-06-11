package abi

import (
	"encoding/hex"

	"github.com/hexoul/eth-rpc-on-aws-lambda/json"
	"github.com/hexoul/eth-rpc-on-aws-lambda/rpc"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

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

func Call(abi abi.ABI, to, name string, inputs []interface{}, outputs interface{}) error {
	data, err := Pack(abi, name, inputs...)
	if err != nil {
		return err
	}

	r := rpc.GetInstance(rpc.Testnet)
	respStr, err := r.Call(to, data)
	if err != nil {
		return err
	}

	resp := json.GetRpcResponseFromJson(respStr)
	return Unpack(abi, outputs, name, resp.Result.(string))
}
