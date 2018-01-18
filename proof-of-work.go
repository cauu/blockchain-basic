package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// 16进制中，每4位对应一个数字，因此此处显示成4的倍数
const targetBits = 16

// 创建一个新的pow运算
func NewProofOfWork(b *Block) *ProofOfWork {
	// 生成一个理论big.Int类型的数字（理论上的可以无限大的int类型）
	target := big.NewInt(1)
	// target向左移动256-targetsBit位
	// 这样我们就创建了一个256位的长整形数字，同时，它的前targetBits都是0
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}

//生成随机数
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	// 将bytes连接到一起
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// 进行pow运算，
// 如果hashInt === pow.target,
// 则表示找到了符合条件的随机数，可以创建区块
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"n", pow.block.HashTransactions)

	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		// compare hashInt和pow.target
		// target的二进制形式形如000010000000000
		// hashInt如果小于target，那么它的前面4位一定是0
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}
