package crypto

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	//_ "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

type Transactions struct {
	hash []common.Hash
}

// Len returns the length of s.
func (s Transactions) Len() int { return len(s.hash) }

// GetRlp implements Rlpable and returns the i'th element of s in rlp.
func (s Transactions) GetRlp(i int) []byte {
	enc, _ := rlp.EncodeToBytes(s.hash[i])
	return enc
}

/*
func DeriveSha(txs []common.Hash) common.Hash {
	transactions := Transactions{hash: txs}
	return types.DeriveSha(transactions)
}
*/

func DeriveSha(txs []common.Hash) (common.Hash, *trie.Trie) {
	transactions := Transactions{hash: txs}
	trie := new(trie.Trie)
	keybuf := new(bytes.Buffer)
	for i := 0; i < transactions.Len(); i++ {
		keybuf.Reset()
		rlp.Encode(keybuf, uint(i))
		trie.Update(keybuf.Bytes(), transactions.GetRlp(i))
		//trie.Update(common.LeftPadBytes([]byte{byte(i)}, 32), transactions.GetRlp(i))
	}
	return trie.Hash(), trie
}

func VerifyProof(txs []common.Hash, tr *trie.Trie) (bool, error) {
	root := tr.Hash()
	proofs := ethdb.NewMemDatabase()
	for _, tx := range txs {
		if tr.Prove(tx.Bytes(), 0, proofs) != nil {
			return false, fmt.Errorf("VerifyProof error missing key %x while constructing proof", tx.Bytes())
		}
		_, _, err := trie.VerifyProof(root, tx.Bytes(), proofs)
		if err != nil {
			return false, fmt.Errorf("VerifyProof error for key %x: %v\nraw proof: %v", tx.Bytes(), err, proofs)
		}
	}
	return true, nil
}
