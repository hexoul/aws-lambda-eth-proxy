package web3

import "testing"

func TestUnit(t *testing.T) {
	units := []string{"gwei", "gweii", "abc", "ether"}
	for _, unit := range units {
		val, err := GetValueOfUnit(unit)
		if val == nil {
			t.Errorf("%s %s", err, unit)
		}
	}
}

func TestFromWei(t *testing.T) {
	FromWei("1234567890000", "ether")
}

func TestToWei(t *testing.T) {
}
