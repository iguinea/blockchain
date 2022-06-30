package blockchain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/iguinea/cryptodemo/pkgs/block"
	"github.com/iguinea/cryptodemo/pkgs/transaction"
	"github.com/iguinea/cryptodemo/pkgs/utils"
)

const (
	MINING_DIFFICULTY = 3
	cMINING_REWARD    = 1.0
	cMINING_SENDER    = "THE BLOCK CHAIN"
	cMINING_TIMER_SEC = 20

	NEIGHBOUR_IP_RANGE_START    = 1
	NEIGHBOUR_IP_RANGE_END      = 10
	BLOCKCHAIN_PORT_RANGE_START = 5000
	BLOCKCHAIN_PORT_RANGE_END   = 5001
)

type Blockchain struct {
	transactionPool []*transaction.Transaction
	//transactionPoolLock *int32

	chain []*block.Block
	//chainLock *int32

	blockChainAddress     string
	localInterfaceAddress net.Addr // Dirección/Interface donde se arranca el server...
	port                  uint16   // Puerto donde se arranca el servidor

	mux         sync.Mutex
	miningTimer *time.Timer

	neighbours    []string
	muxNeighbours sync.Mutex
}

func NewBlockchain(pBlockChainAddress string, pPort uint16) *Blockchain {
	bc := new(Blockchain)

	//bc.chainLock = new(int32)
	//*bc.chainLock = int32(0)

	//bc.transactionPoolLock = new(int32)
	//*bc.transactionPoolLock = int32(0)

	bc.blockChainAddress = pBlockChainAddress
	bc.port = pPort

	// discover local network address
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			//log.Print(fmt.Errorf("localAddresses: %v\n", err.Error()))
			continue
		}
		for _, a := range addrs {
			if i.Name == "eth0" {
				bc.localInterfaceAddress = a
				log.Printf("detected local address: %v %v\n", i.Name, a)
			}
		}
	}

	bc.createBlock()

	return bc
}

func (bc *Blockchain) SetNeighbors() {

	bc.muxNeighbours.Lock()
	utils.FindNeighbors(
		bc.localInterfaceAddress.String(), bc.port,
		NEIGHBOUR_IP_RANGE_START, NEIGHBOUR_IP_RANGE_END,
		BLOCKCHAIN_PORT_RANGE_START, BLOCKCHAIN_PORT_RANGE_END,
	)
	bc.muxNeighbours.Unlock()
}

func (bc *Blockchain) TransactionPool() []*transaction.Transaction {
	return bc.transactionPool
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Bllocks []*block.Block `json:"chains"`
	}{
		Bllocks: bc.chain,
	})
}
func (bc *Blockchain) createBlock() (rBlock *block.Block) {
	//bc.lockTransactionPool()
	//bc.lockChain()

	nonce := bc.proofOfWork()
	previousHash := bc.lastBlock().Hash()
	rBlock = block.NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, rBlock)
	bc.transactionPool = []*transaction.Transaction{}

	//bc.unlockChain()
	//bc.unlockTransactionPool()

	return
}

func (bc *Blockchain) CreateTransaction(sender string, recipient string, value float32, senderPublickey *ecdsa.PublicKey, signature *utils.Signature) bool {
	isTransacted := bc.AddTransaction(sender, recipient, value, senderPublickey, signature)

	// TODO
	// Sync

	return isTransacted
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32, senderPublickey *ecdsa.PublicKey, signature *utils.Signature) bool {

	//bc.lockTransactionPool()
	t := transaction.NewTransaction(nil, nil, sender, recipient, value)

	if sender == cMINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	if bc.verifyTransactionSignature(senderPublickey, signature, t) {
		/*
			if bc.CalculateTotalAmount(sender) < value {
				log.Print("ERROR: Not enough balance in a wallet")
				return false
			}
		*/
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}
	log.Print("ERROR: Verify transaction")
	return false
	//bc.unlockTransactionPool()

}

func (bc *Blockchain) verifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *transaction.Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

func (bc *Blockchain) Mining() bool {
	bc.mux.Lock()
	defer bc.mux.Unlock()

	if len(bc.transactionPool) == 0 {
		return false
	}

	//bc.AddTransaction(cMINING_SENDER, bc.blockChainAddress, cMINING_REWARD)
	bc.AddTransaction(
		cMINING_SENDER,
		bc.blockChainAddress,
		cMINING_REWARD,
		nil,
		nil,
	)
	bc.createBlock()

	old_preffix := log.Prefix()
	log.SetPrefix(fmt.Sprintf("%.16s ", "BlockChain     "))
	log.Print("action=mining, status=success")
	log.SetPrefix(old_preffix)
	return true
}

func (bc *Blockchain) StartMining() {
	bc.Mining()
	bc.miningTimer = time.AfterFunc(time.Second*cMINING_TIMER_SEC, bc.StartMining)
}

func (bc *Blockchain) StoptMining() {
	bc.miningTimer.Stop()
}

func (bc *Blockchain) CalculateTotalAmount(blockchainaddress string) float32 {
	var totalAmount float32 = 0.0

	for _, b := range bc.chain {
		for _, t := range b.GetTransactions() {
			if t.GetRecipient() == blockchainaddress {
				totalAmount += t.GetValue()
			}
			if t.GetSender() == blockchainaddress {
				totalAmount -= t.GetValue()

			}
		}
	}
	return totalAmount
}

func (bc *Blockchain) lastBlock() *block.Block {
	if len(bc.chain) == 0 {
		return &block.Block{}
	}
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) copyTransactionPool() (rTransactions []*transaction.Transaction) {
	rTransactions = make([]*transaction.Transaction, 0)
	for _, t := range bc.transactionPool {
		rTransactions = append(rTransactions, t.CopyTransaction())
	}
	return
}

func (bc *Blockchain) validProof(nonce int, previousHash [32]byte, transactions []*transaction.Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guess := block.AssembleBlock(0, nonce, previousHash, transactions)
	guessHashStr := fmt.Sprintf("%x", guess.Hash())
	return guessHashStr[:difficulty] == zeros
}

func (bc *Blockchain) proofOfWork() int {
	transactions := bc.copyTransactionPool()
	previousHash := bc.lastBlock().Hash()
	nonce := 0
	for !bc.validProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1

	}
	return nonce
}

func (bc *Blockchain) Print() {
	old_preffix := log.Prefix()
	log.SetPrefix(fmt.Sprintf("%.16s ", "BlockChain     "))
	log.Printf("* START PRINTING BLOCKCHAIN\n")

	for i, block := range bc.chain {
		log.Printf("%s Chain %.9d %s\n", strings.Repeat("=", 25), i,
			strings.Repeat("=", 25))
		block.Print()
	}
	log.Printf("* END PRINTING BLOCKCHAIN\n")
	log.Printf("\n")
	log.SetPrefix(old_preffix)

}

type TransactionRequest struct {
	SenderBlockchainAddress    *string  `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string  `json:"recipient_blockchain_address"`
	SenderPublicKey            *string  `json:"sender_public_key"`
	Value                      *float32 `json:"value"`
	Signature                  *string  `json:"signature"`
}

func (tr *TransactionRequest) Validate() bool {
	if tr.Signature == nil ||
		tr.SenderBlockchainAddress == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.SenderPublicKey == nil ||
		tr.Value == nil {
		return false
	}
	return true

}

type AmountResponse struct {
	Amount float32 `json:"amount"`
}

func (ar *AmountResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Amount float32 `json:"amount"`
	}{
		Amount: ar.Amount,
	})
}
