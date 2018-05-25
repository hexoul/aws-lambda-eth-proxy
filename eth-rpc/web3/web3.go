package web3

import (
	"math/big"

	etherCommon "github.com/ethereum/go-ethereum/common"
	common "github.com/hexoul/eth-rpc-on-aws-lambda/eth-rpc/common"
)

func GetValueOfUnit(unit string) (val *big.Int, err string) {
	val = common.UnitMap[unit]
	if val == nil {
		err = "Invalid unit"
	} else if val.Cmp(etherCommon.Big0) == -1 {
		val = nil
		err = "int64 overflow"
	}
	return
}

func FromWei(number int, unit int) int {
	/*
		ret := new(big.Int)
		ret, err := ret.SetString(valStr, 10)
		if !err {
			fmt.Println("Failed to get big from string")
			return nil
		}
	*/
	return 0
}

func ToWei(number int, unit int) int {
	return 0
}
