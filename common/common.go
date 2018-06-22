// Package common includes constants such as unit (wei, ether)
package common

// FindOffsetNBase finds offset and base to parse string to int
func FindOffsetNBase(input string) (offset, base int) {
	if len(input) >= 2 && input[:2] == "0x" {
		return 2, 16
	}
	return 0, 10
}
