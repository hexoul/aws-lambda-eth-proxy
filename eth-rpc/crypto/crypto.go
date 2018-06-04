package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"sync"

	"github.com/hexoul/eth-rpc-on-aws-lambda/eth-rpc/common"
	"github.com/hexoul/eth-rpc-on-aws-lambda/eth-rpc/db"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Crypto struct {
	secretKey string
	nonce     string
	privKey   string
}

// For singleton
var instance *Crypto
var once sync.Once

const (
	DbSecretKeyPropName = "secret_key"
	DbNoncePropName     = "nonce"
	DbPrivKeyPropName   = "priv_key"
)

func GetInstance() *Crypto {
	once.Do(func() {
		dbSecretKey := getPrivKeyFromDB(DbSecretKeyPropName)
		dbNonce := getPrivKeyFromDB(DbNoncePropName)
		dbPrivKey := getPrivKeyFromDB(DbPrivKeyPropName)

		bNonce, _ := hex.DecodeString(dbNonce)
		nPrivKey := DecryptAes(dbPrivKey, dbSecretKey, bNonce)

		fmt.Println(dbSecretKey)
		fmt.Println(dbNonce)
		fmt.Println(nPrivKey)

		instance = &Crypto{
			secretKey: dbSecretKey,
			nonce:     dbNonce,
			privKey:   nPrivKey,
		}
	})
	return instance
}

func getPrivKeyFromDB(propVal string) string {
	//dbHelper := db.GetInstance("aws-region")
	dbHelper := db.GetInstance("")
	if dbHelper == nil {
		return ""
	}

	ret := dbHelper.GetItem(common.DbConfigTblName, common.DbConfigPropName, propVal, common.DbConfigValName)
	item := common.DbConfigResult{}
	for _, elem := range ret.Items {
		dbHelper.UnmarshalMap(elem, &item)
		return item.Value
	}
	return ""
}

func Sign() {
	privKey := getPrivKeyFromDB("priv_key")
	if privKey == "" {
		return
	}
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
	if sig[64] != 27 && sig[64] != 28 {
		return "", fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sig[64] -= 27 // Transform yellow paper V from 27/28 to 0/1

	rpk, err := crypto.Ecrecover(signHash(data), sig)
	if err != nil {
		return "", err
	}
	pubKey := crypto.ToECDSAPub(rpk)
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return fmt.Sprintf("0x%x", recoveredAddr), nil
}

func EcRecoverToPubkey(hash, sig string) ([]byte, error) {
	return crypto.Ecrecover(hexutil.MustDecode(hash), hexutil.MustDecode(sig))
}

func PubkeyToAddress(p []byte) ethcommon.Address {
	return ethcommon.BytesToAddress(crypto.Keccak256(p[1:])[12:])
}

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
