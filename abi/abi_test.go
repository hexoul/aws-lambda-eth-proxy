package abi

import "testing"

var (
	testcontractaddr = "0xc6f1fbb70f850c981591f65f73cd158fb38b6807"

	testabi = "[{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"
)

func TestGetAbiFromJson(t *testing.T) {
	abi, err := GetAbiFromJson(testabi)
	if err != nil {
		t.Errorf("Failed to getAbiFromJson")
	}
	t.Logf("%v", abi)
}

/*
func TestGetAbiFromAddress(t *testing.T) {
	GetAbiFromAddress(testcontractaddr)
}
*/
