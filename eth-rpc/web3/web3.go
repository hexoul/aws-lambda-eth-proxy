package web3

import (
	"fmt"
	"math/big"

	common "github.com/hexoul/eth-rpc-on-aws-lambda/eth-rpc/common"
)

func GetValueOfUnit(unit string) (val *big.Float, err string) {
	val = common.UnitFloatMap[unit]
	if val == nil {
		err = "Invalid unit"
	} else if val.Cmp(common.UnitFloatMap["noether"]) == -1 {
		val = nil
		err = "float64 overflow"
	}
	return
}

func FromWei(number string, unit string) (ret string) {
	// Validate unit
	unitVal, errStr := GetValueOfUnit(unit)
	if errStr != "" {
		fmt.Println("Failed to find unit")
		return
	}

	// Parse number string to big integer
	var err error
	val := new(big.Float)
	if number[:2] == "0x" {
		_, _, err = val.Parse(number[:2], 16)
	} else {
		_, _, err = val.Parse(number, 10)
	}

	if err != nil {
		fmt.Println("Failed to convert number")
		/*
			TODO:
			when number overflow, split to upper and lower bits, convert each bits, combine
		*/
		return
	}

	// Divide number by unit
	val.Quo(val, unitVal)
	return val.String()
}

func ToWei(number string, unit string) (ret string) {
	return
}
