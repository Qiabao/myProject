package lib

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"time"
)

var (
	SYSTEM_ISSUER = crypto.Address("KING_OF_TOKEN")
	num_valid_tx   = 0
	num_invalid_tx = 0
)
type task struct {
	to    crypto.Address
	value int
}

type TokenApp struct {
	types.BaseApplication
	Accounts       map[string]int
	start          time.Time
	height         int64
	wallet         *Wallet
	/*dbs            []dbm.DB
	currentBatch   []*badger.Txn
	num_valid_tx   int
	num_invalid_tx int
	num_query      int64
	start          time.Time
	height         int64
	wallet         *Wallet
	channels       []chan *task
	infoChannels   []chan crypto.Address
	infoSystem     chan crypto.Address
	returnChannels []chan []byte
	returnSystem   chan []byte
	endTime        time.Time
	endTimeQuery   time.Time*/
}
func NewTokenApp() *TokenApp {
	return &TokenApp{Accounts: map[string]int{}, height: 0, wallet: LoadWallet("./wallet")}
}

func NewTokenApp2() *TokenApp {
	return &TokenApp{Accounts: map[string]int{}}
}

func (app *TokenApp) BeginBlock(req types.RequestBeginBlock) (rsp types.ResponseBeginBlock) {
	app.start = time.Now()
	app.height++
	fmt.Println("BeginBlock:", app.start)
	return types.ResponseBeginBlock{}
}

func (app *TokenApp) DeliverTx(req types.RequestDeliverTx) (rsp types.ResponseDeliverTx) {

	return
}
func (app *TokenApp) CheckTx(req types.RequestCheckTx) (rsp types.ResponseCheckTx) {
	/*raw := req.Tx
	tx, err := app.decodeTx(raw)
	if err != nil {
		rsp.Code = 1
		rsp.Log = "decode error"
		return
	}
	if !tx.Verify() {
		rsp.Code = 2
		rsp.Log = "verify failed"
		return
	}*/
	return
}

func (app *TokenApp) decodeTx(raw []byte) (*Tx, error) {
	var tx Tx
	err := codec.UnmarshalJSON(raw, &tx)
	return &tx, err
}


func (app *TokenApp) Query(req types.RequestQuery) (rsp types.ResponseQuery) {
	//??
	addr := crypto.Address(req.Data)
	rsp.Key = req.Data
	rsp.Value, _ = codec.MarshalBinaryBare(app.Accounts[addr.String()])
	return
}


func (app *TokenApp) transfer(from, to crypto.Address, value int) error {
	if app.Accounts[from.String()] < value {
		return errors.New("balance low")
	}
	app.Accounts[from.String()] -= value
	app.Accounts[to.String()] += value
	return nil
}

func (app *TokenApp) issue(issuer, to crypto.Address, value int) error {
	wallet := LoadWallet("./wallet")
	SYSTEM_ISSUER = wallet.GetAddress("issuer")
	if !bytes.Equal(issuer, SYSTEM_ISSUER) {
		return errors.New("invalid issuer")
	}
	app.Accounts[to.String()] += value
	return nil
}

func (app *TokenApp) Dump() {
	fmt.Printf("state => %v\n", app.Accounts)
}
