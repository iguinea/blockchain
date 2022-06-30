package walletserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"text/template"

	"github.com/iguinea/cryptodemo/pkgs/blockchain"
	"github.com/iguinea/cryptodemo/pkgs/transaction"
	"github.com/iguinea/cryptodemo/pkgs/utils"
	"github.com/iguinea/cryptodemo/pkgs/wallet"
)

const tempDir = "pkgs/walletserver/templates/"

type WalletServer struct {
	port    uint16
	gateway string
}

func NewWalletServer(port uint, gateway string) *WalletServer {
	return &WalletServer{port: uint16(port), gateway: gateway}
}

func (ws *WalletServer) Port() uint16    { return ws.port }
func (ws *WalletServer) Gateway() string { return ws.gateway }
func (ws *WalletServer) Index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, err := template.ParseFiles(path.Join(tempDir + "index.html"))
		if err != nil {
			log.Printf("ERROR: template.ParseFiles: %v", err)
			return
		}
		t.Execute(w, "")
	default:
		log.Print("ERROR: invalid HTTP method")
	}
}

func (ws *WalletServer) Wallet(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		myWallet := wallet.NewWallet()
		log.Printf("NewWallet(): %+v", myWallet)
		m, _ := myWallet.MarshallJSON()
		io.WriteString(w, string(m[:]))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Print("ERROR: invalid HTTP method")

	}

}

func (ws *WalletServer) CreateTransaction(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:

		decoder := json.NewDecoder(req.Body)
		var t wallet.TransactionRequest
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
		privateKey := utils.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)
		value, err := strconv.ParseFloat(*t.Value, 32)
		if err != nil {
			log.Print("ERROR: parse error")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		value32 := float32(value)

		w.Header().Add("Content-Type", "application/json")

		transaction := transaction.NewTransaction(privateKey, publicKey, *t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, value32)
		signature := transaction.GenerateSignature()
		signatureStr := signature.String()

		bt := &blockchain.TransactionRequest{
			SenderBlockchainAddress:    t.SenderBlockchainAddress,
			RecipientBlockchainAddress: t.RecipientBlockchainAddress,
			SenderPublicKey:            t.SenderPublicKey,
			Value:                      &value32,
			Signature:                  &signatureStr,
		}
		m, _ := json.Marshal(bt)
		buf := bytes.NewBuffer(m)

		resp, _ := http.Post(ws.Gateway()+"/transactions", "application/json", buf)
		if resp.StatusCode == 201 {
			io.WriteString(w, string(utils.JsonStatus("success")))
			return
		}
		io.WriteString(w, string(utils.JsonStatus("fail")))

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Print("ERROR: invalid HTTP method")

	}
}

func (ws *WalletServer) WalletAmount(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodGet:

		blockchainAdress := req.URL.Query().Get("blockchain_address")
		log.Printf("blockchainAdress: %+v", blockchainAdress)
		endpoint := fmt.Sprintf("%s/amount", ws.gateway)
		//log.Printf("endpoint: %+v", endpoint)

		client := &http.Client{}
		bcsReq, _ := http.NewRequest("GET", endpoint, nil)
		q := bcsReq.URL.Query()
		q.Add("blockchain_address", blockchainAdress)
		bcsReq.URL.RawQuery = q.Encode()
		log.Printf("bcsReq: %+v", bcsReq)

		bcsResp, err := client.Do(bcsReq)

		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		if bcsResp.StatusCode == 200 {
			decoder := json.NewDecoder(bcsResp.Body)
			var bar blockchain.AmountResponse
			err := decoder.Decode(&bar)
			if err != nil {
				log.Printf("ERROR: %v", err)
				io.WriteString(w, string(utils.JsonStatus("fail")))
				return
			}
			log.Printf("bar: %+v", bar)
			m, _ := json.Marshal(struct {
				Message string  `json:"message"`
				Amount  float32 `json:"amount"`
			}{
				Message: "success",
				Amount:  bar.Amount,
			})
			io.WriteString(w, string(m[:]))
		} else {
			io.WriteString(w, string(utils.JsonStatus("fail")))
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Print("ERROR: invalid HTTP method")

	}
}

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/wallet/amount", ws.WalletAmount)
	http.HandleFunc("/transaction", ws.CreateTransaction)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.port)), nil))
}
