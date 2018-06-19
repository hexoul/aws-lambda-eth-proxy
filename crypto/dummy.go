package crypto

import "math/big"

// GetDummy returns dummy Crypto instance for test
func GetDummy() *Crypto {
	return &Crypto{
		secretKey: "dummysecret",
		nonce:     "dummynonce",
		privKey:   "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032",
		Address:   "0x06839e455e0a821f946979d99abe8c4dfdd6fe8b",
		ChainId:   big.NewInt(127),
	}
}
