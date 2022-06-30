package transaction

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/iguinea/cryptodemo/pkgs/utils"
)

type Transaction struct {
	senderPrivateKey           *ecdsa.PrivateKey
	senderPublicKey            *ecdsa.PublicKey
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
}

func NewTransaction(
	privateKey *ecdsa.PrivateKey,
	publicKey *ecdsa.PublicKey,
	sender string,
	recipient string,
	value float32) *Transaction {
	return &Transaction{
		senderPrivateKey:           privateKey,
		senderPublicKey:            publicKey,
		senderBlockchainAddress:    sender,
		recipientBlockchainAddress: recipient,
		value:                      value,
	}
}

/*
func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{
		senderBlockchainAddress:    sender,
		recipientBlockchainAddress: recipient,
		value:                      value,
	}
}
*/
func (t *Transaction) Print(pIndex int) {
	old_preffix := log.Prefix()
	log.SetPrefix(fmt.Sprintf("%.16s ", "Transaction               "))
	log.Printf("  %s tx %.6d %s\n", strings.Repeat("-", 20), pIndex, strings.Repeat("-", 20))
	log.Printf("  sender_blockchain_address     %s\n", t.senderBlockchainAddress)
	log.Printf("  recipient_blockchain_address  %s\n", t.recipientBlockchainAddress)
	log.Printf("  value                         %.3f\n", t.value)
	log.SetPrefix(old_preffix)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}

func (t *Transaction) CopyTransaction() *Transaction {
	return NewTransaction(nil, nil, t.senderBlockchainAddress, t.recipientBlockchainAddress, t.value)

}

func (t *Transaction) GetSender() string {
	return t.senderBlockchainAddress
}

func (t *Transaction) GetRecipient() string {
	return t.recipientBlockchainAddress
}

func (t *Transaction) GetValue() float32 {
	return t.value
}

func (t *Transaction) GenerateSignature() *utils.Signature {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	return &utils.Signature{R: r, S: s}
}

func (t *Transaction) MarshallJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}
