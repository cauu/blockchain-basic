package main

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