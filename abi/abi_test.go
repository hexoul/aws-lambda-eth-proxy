package abi

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

const testabijson = `
[{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"previousOwner","type":"address"},{"indexed":true,"name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"}]
`

const testabijson2 = `
[
	{ "type" : "function", "name" : "balance", "constant" : true },
	{ "type" : "function", "name" : "send", "constant" : false, "inputs" : [ { "name" : "amount", "type" : "uint256" } ] },
	{ "type" : "function", "name" : "test", "constant" : false, "inputs" : [ { "name" : "number", "type" : "uint32" } ] },
	{ "type" : "function", "name" : "string", "constant" : false, "inputs" : [ { "name" : "inputs", "type" : "string" } ] },
	{ "type" : "function", "name" : "bool", "constant" : false, "inputs" : [ { "name" : "inputs", "type" : "bool" } ] },
	{ "type" : "function", "name" : "address", "constant" : false, "inputs" : [ { "name" : "inputs", "type" : "address" } ] },
	{ "type" : "function", "name" : "uint64[2]", "constant" : false, "inputs" : [ { "name" : "inputs", "type" : "uint64[2]" } ] },
	{ "type" : "function", "name" : "uint64[]", "constant" : false, "inputs" : [ { "name" : "inputs", "type" : "uint64[]" } ] },
	{ "type" : "function", "name" : "foo", "constant" : false, "inputs" : [ { "name" : "inputs", "type" : "uint32" } ] },
	{ "type" : "function", "name" : "bar", "constant" : false, "inputs" : [ { "name" : "inputs", "type" : "uint32" }, { "name" : "string", "type" : "uint16" } ] },
	{ "type" : "function", "name" : "slice", "constant" : false, "inputs" : [ { "name" : "inputs", "type" : "uint32[2]" } ] },
	{ "type" : "function", "name" : "slice256", "constant" : false, "inputs" : [ { "name" : "inputs", "type" : "uint256[2]" } ] },
	{ "type" : "function", "name" : "sliceAddress", "constant" : false, "inputs" : [ { "name" : "inputs", "type" : "address[]" } ] },
	{ "type" : "function", "name" : "sliceMultiAddress", "constant" : false, "inputs" : [ { "name" : "a", "type" : "address[]" }, { "name" : "b", "type" : "address[]" } ] }
]`

var (
	testcontractaddr = "0xc6f1fbb70f850c981591f65f73cd158fb38b6807"

	testaddr = "0xd396348325532a21ab2b01aeee1499a713453e7c"
)

func TestGetAbiFromJson(t *testing.T) {
	_, err := GetAbiFromJSON(testabijson2)
	if err != nil {
		t.Fatalf("Failed to GetAbiFromJSON: %s", err)
	}
}

func TestPack(t *testing.T) {
	abi, err := GetAbiFromJSON(testabijson2)
	if err != nil {
		t.Fatalf("Failed to GetAbiFromJSON: %s", err)
	}

	addr := common.HexToAddress(testaddr[2:])
	data, err := Pack(abi, "address", addr)
	if err != nil {
		t.Fatalf("Failed to Pack: %s", err)
	}
	t.Logf("data: %s", data)
}

func TestCall(t *testing.T) {
	abi, err := GetAbiFromJSON(testabijson)
	if err != nil {
		t.Fatalf("Failed to GetAbiFromJSON")
	}

	resp, err := Call(abi, testcontractaddr, "owner", []interface{}{})
	if err != nil || resp.Result == nil || resp.Result == "" || resp.Error == nil || resp.Error.Code != 0 {
		t.Fatalf("Failed to Call %s", err)
	}
	t.Logf("%s", resp.String())

	var addr common.Address
	Unpack(abi, &addr, "owner", resp.Result.(string))
	if len(addr) == 0 {
		t.Fatalf("Failed to Unpack")
	}
}

func TestSendTransaction(t *testing.T) {
	abi, err := GetAbiFromJSON(testabijson)
	if err != nil {
		t.Fatalf("Failed to GetAbiFromJSON")
	}

	addr := common.HexToAddress(testaddr[2:])
	resp, err := DummySendTransaction(abi, testcontractaddr, "transferOwnership", []interface{}{addr}, 0x1)
	if err != nil || resp.Result == "" || resp.Error == nil || resp.Error.Code != 0 {
		t.Fatalf("Failed to SendTransaction %v", resp.Error)
	}
	t.Logf("%s", resp.String())
}

func TestSendTransactionWithSign(t *testing.T) {
	abi, err := GetAbiFromJSON(testabijson)
	if err != nil {
		t.Fatalf("Failed to GetAbiFromJSON")
	}

	addr := common.HexToAddress(testaddr[2:])
	resp, err := DummySendTransactionWithSign(abi, testcontractaddr, "transferOwnership", []interface{}{addr}, 0xffff, 0x1)
	if err != nil || resp.Result == "" || resp.Error == nil || resp.Error.Code != 0 {
		t.Fatalf("Failed to SendTransactionWithSign %v", resp.Error)
	}
	t.Logf("%s", resp.String())
}

/*
func TestGetAbiFromAddress(t *testing.T) {
	GetAbiFromAddress(testcontractaddr)
}
*/
