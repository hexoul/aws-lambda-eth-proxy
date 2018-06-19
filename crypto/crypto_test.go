package crypto

import (
	"bytes"
	"crypto/md5"
	crand "crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	_ "github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

var (
	testmsg     = hexutil.MustDecode("0xce0677bb30baa8cf067c88db9811f4333d131bf8bcf12fe7065d211dce971008")
	testsig     = hexutil.MustDecode("0x90f27b8b488db00b00606796d2987f6a5f59ae62ea05effe84fef5b8b0e549984a691139ad57a3f0b906637673aa2f63d1f55cb1a69199d4009eea23ceaddc9301")
	testpubkey  = hexutil.MustDecode("0x04e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a0a2b2667f7e725ceea70c673093bf67663e0312623c8e091b13cf2c0f11ef652")
	testpubkeyc = hexutil.MustDecode("0x02e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a")

	testmsgraw2 = "0x00000000000000000000000000000000000000000000000000000000deadbeaf"
	testmsg2    = hexutil.MustDecode(testmsgraw2)
	testsigraw2 = "0xbc55faa778761e463be2e06cd3a6ceea4f9113f3f13358776a0484c3bf4d45f938f2496a9e0632edddc0edb9185f9a471e7a6c0f5dea2de8c4b77450e942496d1b"
	testsig2    = hexutil.MustDecode(testsigraw2)
	testpubkey2 = hexutil.MustDecode("0x04fe051b4c866251d356f755dd4d0064da9ff15d24624f336d789448c545845ef2f6d72f98f07a2954c6d77d7844b576e1beae87d60dd7e266f46339f808783e85")

	testmsgraw3 = "0xdeadbeaf"
	testsigraw3 = "0xfb7e213c96e8445737c7fc15cc3674553a4a0c9e4e861e32ad8edbffdae61b1c08aa6bd56db69db045a0778f828e5fb5b41a461fdf2b06a576229784b345eb5b1b"
	testaddr    = "0xd396348325532a21ab2b01aeee1499a713453e7c"

	testprivhex = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	testaddr2   = "0xa9be5a5c3c1862f378409d0fedf2e517a47ac4f2"
)

type kv struct {
	k, v []byte
	t    bool
}

func TestGetPrivKey(t *testing.T) {
}

func TestDeriveShaFromBytes(t *testing.T) {
	var txs []common.Hash
	raws := [][]byte{testsig}
	for _, raw := range raws {
		tx := common.BytesToHash(raw)
		txs = append(txs, tx)
	}
	transactions := Transactions{hash: txs}
	hash := types.DeriveSha(transactions)
	t.Logf("root hash: %x", hash)
}

func TestDeriveShaFromString(t *testing.T) {
	var txs []common.Hash
	raws := []string{"abc", "123"}
	for _, raw := range raws {
		data := crypto.Keccak256([]byte(raw))
		tx := common.BytesToHash(data)
		txs = append(txs, tx)
	}
	//hash := types.DeriveSha(transactions)
	root, tr := DeriveSha(txs)
	t.Logf("root hash: %x", root)
	answer := common.BytesToHash(hexutil.MustDecode("0x1ada6a49c824030f37e8588704a47ee39eab19200d71e8244324ae3f1b146fc9"))
	if root != answer {
		t.Fatalf("Root hash is different")
	}

	if ret, err := VerifyProof(txs, tr); !ret {
		t.Fatalf("Failed to verify proof %s", err)
	}
}

func TestProofTrie(t *testing.T) {
	tr := new(trie.Trie)
	value := &kv{common.LeftPadBytes([]byte{0}, 32), hexutil.MustDecode("0x4e03657aea45a94fc7d47ba826c8d667c0d1e6e33a64a036ec44f58fa12d6c45"), false}
	tr.Update(value.k, value.v)
	value = &kv{common.LeftPadBytes([]byte{1}, 32), hexutil.MustDecode("0x64e604787cbf194841e7b68d7cd28786f6c9a0a3ab9f8b0a0e87cb4387ab0107"), false}
	tr.Update(value.k, value.v)
	t.Logf("%x", tr.Hash())
}

func TestProofRandomTrie(t *testing.T) {
	tr, vals := randomTrie(500)
	root := tr.Hash()
	for _, kv := range vals {
		proofs := ethdb.NewMemDatabase()
		if tr.Prove(kv.k, 0, proofs) != nil {
			t.Fatalf("missing key %x while constructing proof", kv.k)
		}
		val, _, err := trie.VerifyProof(root, kv.k, proofs)
		if err != nil {
			t.Fatalf("VerifyProof error for key %x: %v\nraw proof: %v", kv.k, err, proofs)
		}
		if !bytes.Equal(val, kv.v) {
			t.Fatalf("VerifyProof returned wrong value for key %x: got %x, want %x", kv.k, val, kv.v)
		}
	}
}

func TestSign(t *testing.T) {
	sig, err := Sign(testmsgraw2[2:], testprivhex)
	if err != nil {
		t.Errorf("Failed to sign %s", err)
	}
	addr, ecErr := EcRecover(testmsgraw2, hexutil.Encode(sig))
	if ecErr != nil {
		t.Errorf("Failed to sign, ecrecover error %s", ecErr)
	} else if addr != testaddr2 {
		t.Errorf("Failed to sign, ecrecover mismatch have(%s) want(%s)", addr, testaddr2)
	}
}

func TestAes(t *testing.T) {
	secretKey := "6368616e676520746869732070617373776f726420746f206120736563726574"
	text := "abcde"
	cipher, nonce := EncryptAes(text, secretKey, "")
	ret := DecryptAes(cipher, secretKey, nonce)
	if text != ret {
		t.Errorf("Failed to decrypt")
	}

	cipher, nonce = EncryptAes(text, secretKey, "cd2e39750409adc5f8299c4b")
	ret = DecryptAes(cipher, secretKey, nonce)
	if text != ret {
		t.Errorf("Failed to decrypt")
	}
}

func TestMd5(t *testing.T) {
	text := "abcde"
	hasher := md5.New()
	hasher.Write([]byte(text))
	ret := hex.EncodeToString(hasher.Sum(nil))
	t.Logf("%s", ret)
}

func TestEcRecoverPubkey(t *testing.T) {
	pubkey, err := crypto.Ecrecover(testmsg, testsig)
	if err != nil {
		t.Fatalf("recover error: %s", err)
	}
	if !bytes.Equal(pubkey, testpubkey) {
		t.Errorf("pubkey mismatch: want: %x have: %x", testpubkey, pubkey)
	}

	sig2 := testsig2
	sig2[len(sig2)-1] -= 27
	pubkey, err = crypto.Ecrecover(testmsg2, sig2)
	if err != nil {
		t.Fatalf("recover error: %s", err)
	}
	if !bytes.Equal(pubkey, testpubkey2) {
		t.Errorf("pubkey mismatch: want: %x have: %x", testpubkey2, pubkey)
	}
}

func TestEcRecover(t *testing.T) {
	addr, err := EcRecover(testmsgraw2, testsigraw2)
	if err != nil || addr != testaddr {
		t.Errorf("Failed to EcRecover %s", err)
	}

	addr, err = EcRecover(testmsgraw3, testsigraw3)
	if err != nil || addr != testaddr {
		t.Errorf("Failed to EcRecover %s", err)
	}
}

func TestVerifySignature(t *testing.T) {
	sig := testsig[:len(testsig)-1] // remove recovery id
	if !crypto.VerifySignature(testpubkey, testmsg, sig) {
		t.Errorf("can't verify signature with uncompressed key")
	}
	if !crypto.VerifySignature(testpubkeyc, testmsg, sig) {
		t.Errorf("can't verify signature with compressed key")
	}

	if crypto.VerifySignature(nil, testmsg, sig) {
		t.Errorf("signature valid with no key")
	}
	if crypto.VerifySignature(testpubkey, nil, sig) {
		t.Errorf("signature valid with no message")
	}
	if crypto.VerifySignature(testpubkey, testmsg, nil) {
		t.Errorf("nil signature valid")
	}
	if crypto.VerifySignature(testpubkey, testmsg, append(common.CopyBytes(sig), 1, 2, 3)) {
		t.Errorf("signature valid with extra bytes at the end")
	}
	if crypto.VerifySignature(testpubkey, testmsg, sig[:len(sig)-2]) {
		t.Errorf("signature valid even though it's incomplete")
	}
	wrongkey := common.CopyBytes(testpubkey)
	wrongkey[10]++
	if crypto.VerifySignature(wrongkey, testmsg, sig) {
		t.Errorf("signature valid with with wrong public key")
	}
}

func TestVerifySignature2(t *testing.T) {
	if !crypto.VerifySignature(testpubkey2, testmsg2, testsig2[:len(testsig2)-1]) {
		t.Errorf("invalid signature: pub: %x hash: %x sig: %x", testpubkey2, testmsg2, testsig2)
	}
}

func randomTrie(n int) (*trie.Trie, map[string]*kv) {
	trie := new(trie.Trie)
	vals := make(map[string]*kv)
	/*
		for i := byte(0); i < 100; i++ {
			value := &kv{common.LeftPadBytes([]byte{i}, 32), []byte{i}, false}
			value2 := &kv{common.LeftPadBytes([]byte{i + 10}, 32), []byte{i}, false}
			trie.Update(value.k, value.v)
			trie.Update(value2.k, value2.v)
			vals[string(value.k)] = value
			vals[string(value2.k)] = value2
		}
	*/
	for i := 0; i < n; i++ {
		value := &kv{randBytes(32), randBytes(20), false}
		trie.Update(value.k, value.v)
		vals[string(value.k)] = value
	}
	return trie, vals
}

func randBytes(n int) []byte {
	r := make([]byte, n)
	crand.Read(r)
	return r
}
