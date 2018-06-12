package abi

import (
	"github.com/hexoul/aws-lambda-eth-proxy/crypto"
	"github.com/hexoul/aws-lambda-eth-proxy/json"
	"github.com/hexoul/aws-lambda-eth-proxy/rpc"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func DummySendTransaction(abi abi.ABI, to, name string, inputs []interface{}, gas int) (resp json.RpcResponse, err error) {
	data, err := Pack(abi, name, inputs...)
	if err != nil {
		return
	}

	c := crypto.GetDummy()
	r := rpc.GetInstance(Targetnet)
	respStr, err := r.SendTransaction(c.Address, to, data, gas)
	if err != nil {
		return
	}

	resp = json.GetRpcResponseFromJson(respStr)
	return
}
