package web3

import (
	"fmt"
	"math/big"

	"github.com/hexoul/eth-rpc-on-aws-lambda/eth-rpc/common"
)

func GetValueOfUnit(unit string) *big.Int {
	valStr := common.UnitMap[unit]
	if len(valStr) == 0 {
		return nil
	}

	ret := new(big.Int)
	ret, err := ret.SetString(valStr, 10)
	if !err {
		fmt.Println("Failed to get big from string")
		return nil
	}
	return ret
}

func FromWei(number int, unit int) int {
	return 0
}

func ToWei(number int, unit int) int {
	return 0
}
