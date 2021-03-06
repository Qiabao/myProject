package main

import (
	"encoding/json"
	"fmt"
	kf "github.com/tendermint/tendermint/crypto/secp256k1"
)

type Letter struct {
	Msg       []byte
	Signature []byte
	PubKey    kf.PubKey
}

func main() {
	privSender := kf.GenPrivKey()
	pubSender := privSender.PubKey()

	msg := []byte("some text to send")
	sig, err := privSender.Sign(msg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("send msg => %x\n", msg)
	fmt.Printf("signature => %x\n", sig)

	letter := Letter{msg, sig, pubSender.Bytes()}
	bz, err := json.Marshal(letter)
	if err != nil {
		panic(err)
	}
	fmt.Printf("encoded letter => %x\n", bz)

	var received Letter
	err = json.Unmarshal(bz, &received)
	if err != nil {
		panic(err)
	}
	fmt.Printf("decoded letter => %+v\n", received)

	valid := received.PubKey.VerifySignature(received.Msg, received.Signature)
	fmt.Printf("validated => %t\n", valid)
}
