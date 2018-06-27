package crypto

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

// GetDummy returns dummy Crypto instance for test
func GetDummy() *Crypto {
	privKey, _ := crypto.HexToECDSA("25c317c8d0a63c122073ae52984e8477e7fbc322c93a9457c5579ee6e5a813b3")
	instance := &Crypto{
		secretKey: "dummysecret",
		nonce:     "dummynonce",
		privKey:   privKey,
		ChainID:   big.NewInt(127),
	}
	instance.Sign("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	return instance
}
