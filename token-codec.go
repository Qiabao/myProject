package main

import (
	"myProject/lib"

	"bytes"
	"errors"
	"fmt"
	"github.com/tendermint/tendermint/crypto"
	kf "github.com/tendermint/tendermint/crypto/secp256k1"
)

var (
	issuer3        = kf.GenPrivKey()
	SYSTEM_ISSUER3 = issuer3.PubKey().Address()
)

type TokenApp3 struct {
	Accounts map[string]int
}

func NewTokenApp3() *TokenApp3 {
	return &TokenApp3{Accounts: map[string]int{}}
}

func (app *TokenApp3) transfer(from, to crypto.Address, value int) error {
	if app.Accounts[from.String()] < value {
		return errors.New("balance low")
	}
	app.Accounts[from.String()] -= value
	app.Accounts[to.String()] += value
	return nil
}

func (app *TokenApp3) issue(issuer, to crypto.Address, value int) error {
	if !bytes.Equal(issuer, SYSTEM_ISSUER3) {
		return errors.New("invalid issuer")
	}
	app.Accounts[to.String()] += value
	return nil
}

func (app *TokenApp3) Dump() {
	fmt.Printf("state => %v\n", app.Accounts)
}

func main() {
	app := NewTokenApp3()
	p1 := kf.GenPrivKey()
	p2 := kf.GenPrivKey()

	app.issue(SYSTEM_ISSUER3, p1.PubKey().Address(), 1000)
	app.transfer(p1.PubKey().Address(), p2.PubKey().Address(), 100)
	app.Dump()

	txIssue := lib.NewTx(lib.NewIssuePayload(
		issuer3.PubKey().Address(),
		p1.PubKey().Address(),
		1000))
	txIssue.Sign(issuer3)

	rawtx2, err := lib.MarshalJSON(txIssue)   //这个就是json.Marshal
	rawtx, err := lib.MarshalBinary(txIssue)  //这个就是amino的编码方式，比较省内存

	if err != nil {
		panic(err)
	}
	fmt.Printf("issue tx encoded => %v\n", rawtx)
	fmt.Printf("issue tx encoded => %v\n", rawtx2)

	var txReceived lib.Tx
	err = lib.UnmarshalBinary(rawtx, &txReceived)
	if err != nil {
		panic(err)
	}
	fmt.Printf("issue tx decoded => %+v\n", txReceived)
	fmt.Printf("validated => %t\n", txReceived.Verify())
}
