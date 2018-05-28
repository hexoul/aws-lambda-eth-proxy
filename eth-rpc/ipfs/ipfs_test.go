package ipfs

import "testing"

func TestIpfsNew(t *testing.T) {
	s := New()
	if s == nil {
		t.Errorf("Failed to new IPFS shell")
	}
}
