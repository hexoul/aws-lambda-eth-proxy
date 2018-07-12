// Package web3 is converted golang layer from web3.js
package web3

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/hexoul/aws-lambda-eth-proxy/common"
	"github.com/hexoul/aws-lambda-eth-proxy/log"
)

// GetValueOfUnit returns a value about given unit
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
func getBigFloat(number string) *big.Float {
	var err error
	val := new(big.Float)
	if number[:2] == "0x" {
		_, _, err = val.Parse(number[:2], 16)
	} else {
		_, _, err = val.Parse(number, 10)
	}

	if err != nil {
		log.Error("web3: failed to convert number")
		return nil
	}
	return val
}

// FromWei applys unit to wei
func FromWei(number, unit string) (ret string, err error) {
	// Validate unit
	unitVal, errStr := GetValueOfUnit(unit)
	if errStr != "" {
		err = fmt.Errorf("Failed to find unit")
		return
	}

	// Parse number string to big float
	val := getBigFloat(number)
	if val == nil {
		// TODO: when number overflow, split to upper and lower bits, convert each bits, combine
		err = fmt.Errorf("Number is not appropriate for float64")
		return
	}

	// Divide number by unit
	val.Quo(val, unitVal)
	ret = strings.TrimRight(strings.TrimRight(val.Text('f', 10), "0"), ".")
	return
}

// ToWei gets wei from given value and unit
func ToWei(number, unit string) (ret string, err error) {
	// Validate unit
	unitVal, errStr := GetValueOfUnit(unit)
	if errStr != "" {
		err = fmt.Errorf("Failed to find unit")
		return
	}

	// Parse number string to big float
	val := getBigFloat(number)
	if val == nil {
		// TODO: when number overflow, split to upper and lower bits, convert each bits, combine
		err = fmt.Errorf("Number is not appropriate for float64")
		return
	}

	// Multiply number by unit
	val.Mul(val, unitVal)
	ret = strings.Split(val.Text('f', 10), ".")[0]
	return
}
