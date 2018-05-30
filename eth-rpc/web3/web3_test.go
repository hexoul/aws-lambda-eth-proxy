package web3

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/hexoul/eth-rpc-on-aws-lambda/eth-rpc/common"
)

func TestUnit(t *testing.T) {
	units := []string{"gwei", "gweii", "abc", "ether"}
	for _, unit := range units {
		val, err := GetValueOfUnit(unit)
		if val == nil {
			t.Errorf("%s %s", err, unit)
		}
	}
}

func TestHex(t *testing.T) {
	testHex := "12"

	val := new(big.Float)
	val.Parse(testHex, 16)
	valStr := fmt.Sprintf("%f", val)
	if valStr[:2] != "18" {
		t.Errorf("Failed to parse hex %s", valStr)
	}
}

func TestFromWei(t *testing.T) {
	ret, err := FromWei("1234000000000000000", "ether")
	t.Logf("%s", ret)
	if err != "" || ret == "" || ret[:5] != "1.234" {
		t.Errorf("Failed to FromWei %s", ret)
	}

	ret, err = FromWei("1234000000000000000000000", "ether")
	t.Logf("%s", ret)
	if err != "" || ret == "" || ret[:4] != "1234" {
		t.Errorf("Failed to FromWei %s", ret)
	}
}

func TestToWei(t *testing.T) {
	unit := common.UnitStrMap["ether"]
	ret, err := ToWei("1234", "ether")
	t.Logf("%s", ret)
	if err != "" || ret == "" || (len(unit)+3) != len(ret) || ret[:4] != "1234" {
		t.Errorf("Failed to ToWei %s", ret)
	}
}
