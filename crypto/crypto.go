// Package crypto is combined crypto module for both general(AES, ...) and ethereum(Ecrevoer, sign, ...)
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"sync"

	"github.com/hexoul/aws-lambda-eth-proxy/common"
	"github.com/hexoul/aws-lambda-eth-proxy/db"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// Crypto manager
type Crypto struct {
	secretKey string
	nonce     string
	privKey   string
	Address   string
	ChainID   *big.Int
	Txnonce   uint64
	signer    types.Signer
}

// For singleton
var instance *Crypto
var once sync.Once

const (
	// DbSecretKeyPropName is DB column name about secret key
	DbSecretKeyPropName = "secret_key"
	// DbNoncePropName is DB column name about nonce
	DbNoncePropName = "nonce"
	// DbPrivKeyPropName is DB column name about private key
	DbPrivKeyPropName = "priv_key"
)

// GetInstance returns pointer of Crypto instance
// Because DB operations are needed for Crypto initiation,
// Crypto is designed as singleton to reduce the number of DB operation units used
func GetInstance() *Crypto {
	once.Do(func() {
		dbSecretKey := getConfigFromDB(DbSecretKeyPropName)
		dbNonce := getConfigFromDB(DbNoncePropName)
		dbPrivKey := getConfigFromDB(DbPrivKeyPropName)

		var nPrivKey string
		if dbSecretKey != "" && dbNonce != "" && dbPrivKey != "" {
			bNonce, _ := hex.DecodeString(dbNonce)
			nPrivKey = DecryptAes(dbPrivKey, dbSecretKey, bNonce)
		}

		instance = &Crypto{
			secretKey: dbSecretKey,
			nonce:     dbNonce,
			privKey:   nPrivKey,
			Txnonce:   0,
		}
		instance.Sign("0xabcdef")
	})
	return instance
}

// Sign returns signed message using own private key
func (c *Crypto) Sign(msg string) string {
	sig, err := Sign(msg, c.privKey)
	if err != nil {
		return ""
	}
	sig[64] += 27

	ret := hexutil.Encode(sig)
	if c.Address == "" {
		c.Address, _ = EcRecover(msg, ret)
		fmt.Printf("Crypto address is set to %s\n", c.Address)
	}
	return ret
}

// SignTx returns signed transaction using own private key
func (c *Crypto) SignTx(tx *types.Transaction) (*types.Transaction, error) {
	if c.signer != nil {
		// Nothing to do
	} else if c.ChainID != nil {
		c.signer = types.NewEIP155Signer(c.ChainID)
	} else {
		c.signer = types.HomesteadSigner{}
	}
	privKey, _ := crypto.HexToECDSA(c.privKey)
	signedTx, err := types.SignTx(tx, c.signer, privKey)
	if err != nil {
		return nil, fmt.Errorf("tx or private key is not appropriate")
	}
	return signedTx, nil
}

// Sign returns signed message using given private key
func Sign(msg, privKey string) ([]byte, error) {
	key, _ := crypto.HexToECDSA(privKey)
	bMsg := crypto.Keccak256([]byte(msg))
	return crypto.Sign(signHash(bMsg), key)
}

// getConfigFromDB returns value string matching given key at config table
// Basically, connect to DynamoDB placed in same region of lambda executing this function
// To use specific region, change the parameter of db.GetInstance()
func getConfigFromDB(propVal string) string {
	//dbHelper := db.GetInstance("aws-region")
	dbHelper := db.GetInstance("")
	if dbHelper == nil {
		return ""
	}

	ret := dbHelper.GetItem(common.DbConfigTblName, common.DbConfigPropName, propVal, common.DbConfigValName)
	if ret == nil {
		return ""
	}

	item := common.DbConfigResult{}
	for _, elem := range ret.Items {
		dbHelper.UnmarshalMap(elem, &item)
		return item.Value
	}
	return ""
}

// signHash is a helper function that calculates a hash for the given message that can be
// safely used to calculate a signature from.
//
// The hash is calulcated as
//   keccak256("\x19Ethereum Signed Message:\n"${message length}${message}).
//
// This gives context to the signed message and prevents signing of transactions.
func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

// EcRecover returns the address for the Account that was used to create the signature.
// Note, this function is compatible with eth_sign and personal_sign. As such it recovers
// the address of:
// hash = keccak256("\x19Ethereum Signed Message:\n"${message length}${message})
// addr = ecrecover(hash, signature)
//
// Note, the signature must conform to the secp256k1 curve R, S and V values, where
// the V value must be be 27 or 28 for legacy reasons.
//
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_ecRecover
func EcRecover(dataStr, sigStr string) (string, error) {
	data := hexutil.MustDecode(dataStr)
	sig := hexutil.MustDecode(sigStr)
	if len(sig) != 65 {
		return "", fmt.Errorf("signature must be 65 bytes long")
	}
	if sig[64] == 0 || sig[64] == 1 {
		// Nothing to do
	} else if sig[64] == 27 || sig[64] == 28 {
		sig[64] -= 27 // Transform yellow paper V from 27/28 to 0/1
	} else {
		return "", fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}

	rpk, err := crypto.Ecrecover(signHash(data), sig)
	if err != nil {
		return "", err
	}
	pubKey, err := crypto.UnmarshalPubkey(rpk)
	if err != nil {
		return "", err
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return fmt.Sprintf("0x%x", recoveredAddr), nil
}

// EcRecoverToPubkey returns public key through EcRecover
func EcRecoverToPubkey(hash, sig string) ([]byte, error) {
	return crypto.Ecrecover(hexutil.MustDecode(hash), hexutil.MustDecode(sig))
}

// PubkeyToAddress converts public key to ethereum address
func PubkeyToAddress(p []byte) ethcommon.Address {
	return ethcommon.BytesToAddress(crypto.Keccak256(p[1:])[12:])
}

// EncryptAes encrypts text using AES with given key and nonce
func EncryptAes(text, keyStr, nonceStr string) (string, []byte) {
	// Load your secret key from a safe place and reuse it across multiple
	// Seal/Open calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).
	key, _ := hex.DecodeString(keyStr)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	var nonce []byte
	if nonceStr == "" {
		nonce = make([]byte, 12)
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			panic(err.Error())
		}
	} else {
		nonce, _ = hex.DecodeString(nonceStr)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nonce
}

// DecryptAes decrypts text using AES with given key and nonce
func DecryptAes(text, keyStr string, nonce []byte) string {
	key, _ := hex.DecodeString(keyStr)
	ciphertext, _ := hex.DecodeString(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return string(plaintext[:])
}
