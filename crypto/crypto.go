// Package crypto is combined crypto module for both general(AES, ...) and ethereum(Ecrevoer, sign, ...)
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hexoul/aws-lambda-eth-proxy/common"
	"github.com/hexoul/aws-lambda-eth-proxy/db"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// Crypto manager
// It should be initialized first at main
type Crypto struct {
	privKey *ecdsa.PrivateKey
	address string

	signer  types.Signer
	chainID *big.Int

	txnonce uint64
}

// For singleton
var (
	instance       *Crypto
	once           sync.Once
	mutex          = &sync.Mutex{}
	PathChan       = make(chan string)
	PassphraseChan = make(chan string)
)

// For DB columns
const (
	// DbSecretKeyPropName is DB column name about secret key
	DbSecretKeyPropName = "secret_key"
	// DbNoncePropName is DB column name about nonce
	DbNoncePropName = "nonce"
	// DbKeyJSONPropName is DB column name about key json
	DbKeyJSONPropName = "key_json"
)

// For environment arguments
const (
	// Passphrase means passphrase used to decrypt keystore
	Passphrase = "KEY_PASSPHRASE"
	// Path means a location of keyjson in file system
	Path = "KEY_PATH"
	// IsAwsLambda decides if served as AWS lambda or not
	IsAwsLambda = "AWS_LAMBDA"
)

// GetInstance returns pointer of Crypto instance
// Because DB operations are needed for Crypto initiation,
// Crypto is designed as singleton to reduce the number of DB operation units used
func GetInstance() *Crypto {
	// Check if already assigned
	if instance != nil {
		return instance
	}

	// Check channel within the timeout
	var path string
	select {
	case path = <-PathChan:
		break
	case <-time.After(1 * time.Second):
		break
	}
	if path == "" {
		return instance
	}

	// Initalize
	once.Do(func() {
		passphrase := <-PassphraseChan

		var privkey *ecdsa.PrivateKey
		var addr string
		if os.Getenv(IsAwsLambda) != "" {
			privkey, addr = getPrivateKeyFromDB(passphrase)
		} else {
			privkey, addr = getPrivateKeyFromFile(path, passphrase)
		}

		if addr == "" {
			panic("Failed to parse key json for Crypto")
		}

		//fmt.Printf("privkey %s, addr: %s\n", hex.EncodeToString(crypto.FromECDSA(privkey)), addr)
		fmt.Println("Crypto address is set to ", addr)
		instance = &Crypto{
			privKey: privkey,
			address: addr,
		}
	})
	return instance
}

// getPrivateKeyFromDB returns private key and address from DB
func getPrivateKeyFromDB(passphrase string) (privkey *ecdsa.PrivateKey, addr string) {
	dbSecretKey := getConfigFromDB(DbSecretKeyPropName)
	dbNonce := getConfigFromDB(DbNoncePropName)
	dbKeyJSON := getConfigFromDB(DbKeyJSONPropName)
	if dbSecretKey == "" || dbNonce == "" || dbKeyJSON == "" {
		return
	}

	bNonce, _ := hex.DecodeString(dbNonce)
	keyjson := DecryptAes(dbKeyJSON, dbSecretKey, bNonce)
	if key, err := keystore.DecryptKey([]byte(keyjson), passphrase); err == nil {
		return key.PrivateKey, key.Address.String()
	}
	return
}

// getPrivateKeyFromFile returns private key and address from file
func getPrivateKeyFromFile(filepath, passphrase string) (privkey *ecdsa.PrivateKey, addr string) {
	keyjson, err := ioutil.ReadFile(filepath)
	if err != nil {
		return
	}
	if key, err := keystore.DecryptKey(keyjson, passphrase); err == nil {
		return key.PrivateKey, key.Address.String()
	}
	return
}

// InitChainID initalizes chain ID
func (c *Crypto) InitChainID(chainID *big.Int) {
	if c.chainID == nil {
		c.chainID = chainID
	}
}

// InitNonce initailizes TX nonce one time
func (c *Crypto) InitNonce(nonce uint64) {
	if c.txnonce == 0 {
		c.txnonce = nonce
	}
}

// GetAddress returns an address of Crypto manager
func (c *Crypto) GetAddress() string {
	return c.address
}

// Sign returns signed message using own private key
func (c *Crypto) Sign(msg string) string {
	sig, err := Sign(msg, c.privKey)
	if err != nil {
		return ""
	}
	sig[64] += 27

	ret := hexutil.Encode(sig)
	if c.address == "" {
		addr, err := EcRecover(hexutil.Encode(crypto.Keccak256([]byte(msg))), ret)
		if err != nil {
			c.address = fmt.Sprintf("0x%x", addr)
			fmt.Printf("Crypto address is set to %s\n", c.address)
		}
	}
	return ret
}

// SignTx returns signed transaction using own private key
func (c *Crypto) SignTx(tx *types.Transaction) (*types.Transaction, error) {
	if c.signer != nil {
		// Nothing to do
	} else if c.chainID != nil {
		c.signer = types.NewEIP155Signer(c.chainID)
	} else {
		c.signer = types.HomesteadSigner{}
	}
	signedTx, err := types.SignTx(tx, c.signer, c.privKey)
	if err != nil {
		return nil, fmt.Errorf("tx or private key is not appropriate")
	}
	return signedTx, nil
}

// ApplyNonce applies nonce to a given function "f"
// Function description should be func(uint64) (error)
// If given function returns nil error, increase nonce
// Meaning of this function's return is either nonce was increased or not
func (c *Crypto) ApplyNonce(f interface{}) bool {
	mutex.Lock()
	defer mutex.Unlock()
	nonce := atomic.LoadUint64(&c.txnonce)
	err := f.(func(uint64) error)(nonce)
	if err != nil {
		return false
	}
	atomic.AddUint64(&c.txnonce, 1)
	return true
}

// Sign returns signed message using given private key
func Sign(msg string, privKey *ecdsa.PrivateKey) ([]byte, error) {
	bMsg := crypto.Keccak256([]byte(msg))
	return crypto.Sign(signHash(bMsg), privKey)
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
func EcRecover(dataStr, sigStr string) (addr ethcommon.Address, err error) {
	data := hexutil.MustDecode(dataStr)
	sig := hexutil.MustDecode(sigStr)
	if len(sig) != 65 {
		err = fmt.Errorf("signature must be 65 bytes long")
		return
	}
	if sig[64] == 0 || sig[64] == 1 {
		// Nothing to do
	} else if sig[64] == 27 || sig[64] == 28 {
		sig[64] -= 27 // Transform yellow paper V from 27/28 to 0/1
	} else {
		err = fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
		return
	}

	rpk, err := crypto.Ecrecover(signHash(data), sig)
	if err != nil {
		return
	}
	pubKey, err := crypto.UnmarshalPubkey(rpk)
	if err != nil {
		return
	}
	return crypto.PubkeyToAddress(*pubKey), nil
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
