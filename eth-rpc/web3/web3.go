package web3

import (
	"fmt"
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

func FromWei(number string, unit string) (ret string) {
	var err bool
	val := new(big.Int)
	if number[:2] == "0x" {
		_, err = val.SetString(number, 16)
	} else {
		_, err = val.SetString(number, 10)
	}

	if !err {
		fmt.Println("Failed to convert number")
	}
	return
}

func ToWei(number string, unit string) (ret string) {
	return
}
