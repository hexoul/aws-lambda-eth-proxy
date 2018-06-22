package crypto

import (
	"math/big"
)

// GetDummy returns dummy Crypto instance for test
func GetDummy() *Crypto {
	instance := &Crypto{
		secretKey: "dummysecret",
		nonce:     "dummynonce",
		privKey:   "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032",
		ChainID:   big.NewInt(127),
	}
	instance.Sign("0xabcdef")
	return instance
}
