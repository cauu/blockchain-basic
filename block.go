package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"strconv"
	"time"
)

// Block
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// 区块链的两大关键技术
// 1. 为保证区块链数据的可信度，需要通过“共识机制”来保证写入区块链的数据是被大多数人所认可的
// 2. 向区块链中添加新的区块并不是那么容易：用户需要付出一定的努力才能实现
func (b *Block) SetHash() {
	// 拿到当前的timestamp
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	// 将当前的数据和前一个区块的hash拼在一起并计算出当前区块的哈希值
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.HashTransactions(), timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func (b *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)

	return result.Bytes(), err
}

func DeserializeBlock(d []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)

	return &block, err
}

// 在区块链的最后一个区块之后添加一个区块
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	// 如果希望添加block，需要挖矿（证明自己的工作量）
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	// block.SetHash()
	return block
}
