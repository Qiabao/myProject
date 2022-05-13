package lib

import (
	"bytes"
	"github.com/tendermint/tendermint/crypto"
	"time"
)

type Tx struct {
	Payload   Payload            //支付行为
	Signature []byte             //签名
	PubKey    crypto.PubKey      //公钥
	Sequence  int64              //时间
}

//生成交易Tx
func NewTx(payload Payload) *Tx {
	return &Tx{Payload: payload, Sequence: time.Now().Unix()}
}

func (tx *Tx) Verify() bool {
	signer := tx.Payload.GetSigner()
	signerFromKey := tx.PubKey.Address()
	if !bytes.Equal(signer, signerFromKey) {
		return false
	}
	data := tx.Payload.GetSignBytes()
	sig := tx.Signature
	valid := tx.PubKey.VerifySignature(data, sig)
	if !valid {
		return false
	}
	return true
}

func (tx *Tx) Sign(priv crypto.PrivKey) error {
	data := tx.Payload.GetSignBytes()
	var err error
	tx.Signature, err = priv.Sign(data)
	tx.PubKey = priv.PubKey()
	return err
}
