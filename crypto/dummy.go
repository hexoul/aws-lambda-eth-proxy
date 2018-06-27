package crypto

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

// GetDummy returns dummy Crypto instance for test
func GetDummy() *Crypto {
	privKey, _ := crypto.HexToECDSA("25c317c8d0a63c122073ae52984e8477e7fbc322c93a9457c5579ee6e5a813b3")
	instance := &Crypto{
		privKey: privKey,
		ChainID: big.NewInt(127),
		Address: "0xed56062123b0301a9a642f85f2711581bec8d79d",
	}
	return instance
}
