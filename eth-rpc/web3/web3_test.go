package web3

import "testing"

func TestUnit(t *testing.T) {
	units := []string{"gwei", "gweii", "abc", "ether"}
	for _, unit := range units {
		val := GetValueOfUnit(unit)
		if val == nil {
			t.Errorf("There is no unit %s", unit)
		}
	}
}

func TestFromWei(t *testing.T) {
}

func TestToWei(t *testing.T) {
}
