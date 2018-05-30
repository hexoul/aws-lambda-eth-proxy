package web3

import (
	"fmt"
	"math/big"
	"strings"

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

// Parse number string to big float
func GetBigFloat(number string) *big.Float {
	var err error
	val := new(big.Float)
	if number[:2] == "0x" {
		_, _, err = val.Parse(number[:2], 16)
	} else {
		_, _, err = val.Parse(number, 10)
	}

	if err != nil {
		fmt.Println("Failed to convert number")
		return nil
	}
	return val
}

func FromWei(number, unit string) (ret, err string) {
	// Validate unit
	unitVal, errStr := GetValueOfUnit(unit)
	if errStr != "" {
		err = "Failed to find unit"
		fmt.Println(err)
		return
	}

	// Parse number string to big float
	val := GetBigFloat(number)
	if val == nil {
		// TODO: when number overflow, split to upper and lower bits, convert each bits, combine
		err = "Number is not appropriate for float64"
		return
	}

	// Divide number by unit
	val.Quo(val, unitVal)
	ret = strings.TrimRight(strings.TrimRight(val.Text('f', 10), "0"), ".")
	return
}

func ToWei(number, unit string) (ret, err string) {
	// Validate unit
	unitVal, errStr := GetValueOfUnit(unit)
	if errStr != "" {
		err = "Failed to find unit"
		fmt.Println(err)
		return
	}

	// Parse number string to big float
	val := GetBigFloat(number)
	if val == nil {
		// TODO: when number overflow, split to upper and lower bits, convert each bits, combine
		err = "Number is not appropriate for float64"
		return
	}

	// Multiply number by unit
	val.Mul(val, unitVal)
	ret = strings.Split(val.Text('f', 10), ".")[0]
	return
}
