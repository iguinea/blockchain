package block

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/iguinea/cryptodemo/pkgs/transaction"
)

type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*transaction.Transaction
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*transaction.Transaction) *Block {
	return AssembleBlock(time.Now().UnixNano(), nonce, previousHash, transactions)
}

func AssembleBlock(timestamp int64, nonce int, previousHash [32]byte, transactions []*transaction.Transaction) *Block {
	b := new(Block)
	b.timestamp = timestamp
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions
	return b
}
func (b *Block) Print() {
	old_preffix := log.Prefix()
	log.SetPrefix(fmt.Sprintf("%.16s ", "Block               "))
	log.Printf("timestamp       %d\n", b.timestamp)
	log.Printf("nonce           %d\n", b.nonce)
	log.Printf("previous_hash   %x\n", b.previousHash)
	for i, t := range b.transactions {
		t.Print(i + 1)
	}
	//log.Printf("transactions    %s\n", b.transactions)
	log.SetPrefix(old_preffix)
}

func (b *Block) Hash() [32]byte {
	m, _ := b.MarshalJSON()
	return sha256.Sum256([]byte(m))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64                      `json:"timestamp"`
		Nonce        int                        `json:"nonce"`
		PreviousHash string                     `json:"previous_hash"`
		Transactions []*transaction.Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: fmt.Sprintf("%x", b.previousHash),
		Transactions: b.transactions,
	})
}
func (b *Block) GetTransactions() []*transaction.Transaction {
	return b.transactions
}
