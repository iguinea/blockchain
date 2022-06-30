package blockchainserver

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/iguinea/cryptodemo/pkgs/blockchain"
	"github.com/iguinea/cryptodemo/pkgs/transaction"
	"github.com/iguinea/cryptodemo/pkgs/utils"
	"github.com/iguinea/cryptodemo/pkgs/wallet"
)

var cache map[string]*blockchain.Blockchain = make(map[string]*blockchain.Blockchain)

type BlockChainServer struct {
	port uint16
}

func NewBlockChainServer(port uint16) *BlockChainServer {
	return &BlockChainServer{port: port}

}

func (bcs *BlockChainServer) Port() uint16 {
	return bcs.port
}

func (bcs *BlockChainServer) GetBlockChain() *blockchain.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		minersWallet := wallet.NewWallet()
		bc = blockchain.NewBlockchain(minersWallet.BlockChainAddress(), bcs.Port())
		cache["blockchain"] = bc
		log.Printf("[GetBlockChain()] private_key        %v", minersWallet.PrivateKeyStr())
		log.Printf("[GetBlockChain()] public_key         %v", minersWallet.PublicKeyStr())
		log.Printf("[GetBlockChain()] blockchain_address %v", minersWallet.BlockChainAddress())
	}
	return bc
}

// GetChain will return JSON of the blockchain
func (bcs *BlockChainServer) GetChain(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockChain()
		m, _ := bc.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		log.Printf("ERROR: invalid HTTP method: %s", req.Method)
	}

}

// Adds a new transaction
func (bcs *BlockChainServer) Transactions(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		//TODO
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockChain()
		transactions := bc.TransactionPool()
		m, _ := json.Marshal(struct {
			Transactions []*transaction.Transaction `json:"transactions"`
			Length       int                        `json:"length"`
		}{
			Transactions: transactions,
			Length:       len(transactions),
		})
		io.WriteString(w, string(m[:]))
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var t blockchain.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		if !t.Validate() {
			log.Print("ERROR: missing field(s)")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		signature := utils.SignatureFromString(*t.Signature)

		bc := bcs.GetBlockChain()
		isCreated := bc.CreateTransaction(*t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, *t.Value, publicKey, signature)

		w.Header().Add("Content/Type", "application/json")
		var m []byte
		if !isCreated {
			log.Printf("ERROR: transaction is not created")

			w.WriteHeader(http.StatusBadRequest)
			m = utils.JsonStatus("fail")
		} else {
			w.WriteHeader(http.StatusCreated)
			m = utils.JsonStatus("success")
		}
		io.WriteString(w, string(m))

	default:

		log.Printf("ERROR: invalid HTTP method: %s", req.Method)
	}

}

func (bcs *BlockChainServer) Mine(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		bc := bcs.GetBlockChain()
		isMined := bc.Mining()
		var m []byte
		if !isMined {
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JsonStatus("fail")
		} else {
			m = utils.JsonStatus("success")
		}
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m))

	default:
		log.Printf("ERROR: invalid HTTP method: %s", req.Method)
		w.WriteHeader(http.StatusBadRequest)

	}
}

func (bcs *BlockChainServer) StartMining(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		bcs.GetBlockChain().StartMining()
		m := utils.JsonStatus("success")
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m))

	default:
		log.Printf("ERROR: invalid HTTP method: %s", req.Method)
		w.WriteHeader(http.StatusBadRequest)

	}
}

func (bcs *BlockChainServer) StopMining(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		bcs.GetBlockChain().StoptMining()
		m := utils.JsonStatus("success")
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m))

	default:
		log.Printf("ERROR: invalid HTTP method: %s", req.Method)
		w.WriteHeader(http.StatusBadRequest)

	}
}

func (bcs *BlockChainServer) Amount(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		blockchainAddress := req.URL.Query().Get("blockchain_address")
		amount := bcs.GetBlockChain().CalculateTotalAmount(blockchainAddress)

		ar := &blockchain.AmountResponse{Amount: amount}
		m, _ := ar.MarshalJSON()

		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m[:]))

	default:
		log.Printf("ERROR: invalid HTTP method: %s", req.Method)
		w.WriteHeader(http.StatusBadRequest)

	}
}
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "resources/ico.png")
}

func (bcs *BlockChainServer) Run() {
	http.HandleFunc("/", bcs.GetChain)
	http.HandleFunc("/transactions", bcs.Transactions)
	http.HandleFunc("/mine", bcs.Mine)
	http.HandleFunc("/mine/start", bcs.StartMining)
	http.HandleFunc("/mine/stop", bcs.StopMining)
	http.HandleFunc("/amount", bcs.Amount)
	http.HandleFunc("/favicon.ico", faviconHandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcs.port)), nil))
}
