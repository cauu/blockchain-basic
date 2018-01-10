package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

// 区块链的两大关键技术
// 1. 为保证区块链数据的可信度，需要通过“共识机制”来保证写入区块链的数据是被大多数人所认可的
// 2. 向区块链中添加新的区块并不是那么容易：用户需要付出一定的努力才能实现

func (b *Block) SetHash() {
	// 拿到当前的timestamp
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	// 将当前的数据和前一个区块的hash拼在一起并计算出当前区块的哈希值
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

// 在区块链的最后一个区块之后添加一个区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	block.SetHash()
	return block
}

type Blockchain struct {
	blocks []*Block
}

// 新建一个区块(data是需要写入区块的数据)
// 真正的区块链中，此处需要大量的算力来计算
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

// 新建创世区块
func NewGenesisBlock() *Block {
	return NewBlock("Genisis Block", []byte{})
}

//初始化区块链
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func main() {
	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
