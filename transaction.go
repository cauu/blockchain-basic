package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

var subsidy = 10

// 我们之前在比特币中创建的区块，
// 就是用来记录这些交易(transation)而存在的
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

type TXOutput struct {
	// Output中的value是不可分的
	// 例如:
	// 如果当前output中有10个btc, 你希望转5个给别人,
	// 那么你首先会将10个btc都转给对方，同时会生成一个changeTx,
	// 将剩余的5个btc转给你
	Value int
	// ScriptPubKey = OP_DUP OP_HASH160 <pubKeyHash> OP_EQUALVERIFY OP_CHECKSIG
	ScriptPubKey string
}

type TXInput struct {
	Txid      []byte // input对应的output所属的transactionId
	Vout      int    // 前一个transaction中output的index
	ScriptSig string // ScriptSig = <sig><pubKey // 其中pubKey对应发送方A的钱包的公钥
}

// Serialize returns a serialized Transaction
func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte
	txCopy := tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

// 通过挖矿可以获得奖励
func NewCoinBaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.ID = tx.Hash()

	return &tx
}
