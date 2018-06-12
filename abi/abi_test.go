package abi

import "testing"

var (
	testcontractaddr = "0xc6f1fbb70f850c981591f65f73cd158fb38b6807"

	testabi = "[{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

	testaddr = "0xd396348325532a21ab2b01aeee1499a713453e7c"
)

func TestGetAbiFromJson(t *testing.T) {
	abi, err := GetAbiFromJson(testabi)
	if err != nil {
		t.Errorf("Failed to GetAbiFromJson: %s", err)
	}
	t.Logf("%v", abi)
}

func TestPack(t *testing.T) {
	abi, err := GetAbiFromJson(testabi)
	if err != nil {
		t.Errorf("Failed to GetAbiFromJson: %s", err)
	}

	data, err := Pack(abi, "transferOwnership", []interface{}{testaddr})
	if err != nil {
		t.Errorf("Failed to Pack: %s", err)
	}
	t.Logf("%s", data)
}

func TestSendTransaction(t *testing.T) {
	abi, err := GetAbiFromJson(testabi)
	if err != nil {
		t.Errorf("Failed to GetAbiFromJson")
	}

	resp, err := SendTransaction(abi, testcontractaddr, "newOwner", []interface{}{"a", "b"}, 0x1)
	if resp.Result == "" || resp.Error.Code != 0 {
		t.Errorf("Failed to SendTransaction")
	}
}

func TestSendRawTransaction(t *testing.T) {
	_, err := GetAbiFromJson(testabi)
	if err != nil {
		t.Errorf("Failed to GetAbiFromJson")
	}
}

/*
func TestGetAbiFromAddress(t *testing.T) {
	GetAbiFromAddress(testcontractaddr)
}
*/
