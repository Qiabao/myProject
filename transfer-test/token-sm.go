package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tendermint/tendermint/crypto"
)

var (
	SYSTEM_ISSUER2 = crypto.Address("KING_OF_TOKEN")
)

type TokenApp2 struct {
	Accounts map[string]int
}

func NewTokenApp() *TokenApp2 {
	return &TokenApp2{Accounts: map[string]int{}}
}

func (app *TokenApp2) transfer(from, to crypto.Address, value int) error {
	if app.Accounts[from.String()] < value {
		return errors.New("balance low")
	}
	app.Accounts[from.String()] -= value
	app.Accounts[to.String()] += value
	return nil
}

func (app *TokenApp2) issue(issuer, to crypto.Address, value int) error {
	if !bytes.Equal(issuer, SYSTEM_ISSUER2) {
		return errors.New("invalid issuer")
	}
	app.Accounts[to.String()] += value
	return nil
}

func (app *TokenApp2) Dump() {
	fmt.Printf("state => %v\n", app.Accounts)
}

func main() {
	app := NewTokenApp()
	a1 := crypto.Address("TEST_ADDR1")
	a2 := crypto.Address("TEST_ADDR2")

	app.issue(SYSTEM_ISSUER2, a1, 1000)
	app.transfer(a1, a2, 100)
	app.Dump()
}
