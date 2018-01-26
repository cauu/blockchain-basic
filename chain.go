package main

import (
	"encoding/hex"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

var dbFile = "db"
var blocksBucket = "BlockBucket"

const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

// type Blockchain struct {
// 	blocks []*Block
// }

// 此处不再将block chain存放在内存中
// tip -> 最后一个block的hash
// db用于读取数据库中存放的blockChain
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// 1. 新建一个区块(data是需要写入区块的数据)
// 真正的区块链中，此处需要大量的算力来计算
// 2. 区块链的数据被储存在db中，添加新的区块时，
// 首先从数据库中读取最后一个区块的hash
// 再创建新的区块，如果合法，就写入到数据库中
// 3. 将区块信息进行同步
func (bc *Blockchain) AddBlock(data string) error {
	var lastHash []byte
	var err error

	defer func() {
		if err != nil {
			log.Panic(err)
		}
	}()

	err = bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("1"))

		return nil
	})

	// 添加新块时,需要
	cbtx := NewCoinBaseTx("-1", data)
	newBlock := NewBlock([]*Transaction{cbtx}, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		serialized, err := newBlock.Serialize()
		err = b.Put(newBlock.Hash, serialized)
		err = b.Put([]byte("1"), newBlock.Hash)
		bc.tip = newBlock.Hash

		if err != nil {
			log.Panic(err)
		}

		return nil
	})
	// prevBlock := bc.blocks[len(bc.blocks)-1]
	// newBlock := NewBlock(data, prevBlock.Hash)
	// bc.blocks = append(bc.blocks, newBlock)

	return err
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

// 找到所有的未被任何input引用的output transactions
func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction
	// spentTXOs形如[txid: [out1,out2,out3]]
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		// 此处next是从最后一个区块反向遍历至第一个区块
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			// 遍历所有output transactions，找到没被消耗的output
		Outputs:
			for outIdx, out := range tx.Vout {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			// 找到所有被消耗的output
			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}

func (bc *Blockchain) FIndUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// 新建创世区块
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

func dbExist() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

//初始化区块链
func NewBlockchain(address string) *Blockchain {
	// if dbExist() {
	// 	fmt.Println("blockchain already exist.")
	// 	os.Exit(1)
	// }
	var tip []byte
	// 打开dbfile
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		log.Panic(err)
	}

	// 启动一个boltDB的read-write事务
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			cbtx := NewCoinBaseTx(address, genesisCoinbaseData)
			genesis := NewGenesisBlock(cbtx)

			b, err := tx.CreateBucket([]byte(blocksBucket))

			if err != nil {
				log.Panic(err)
			}
			// bucket存放的内容为:
			// block hash -> block serialized
			// 1 -> hash of last block in the chain
			serialized, err := genesis.Serialize()
			err = b.Put(genesis.Hash, serialized)
			err = b.Put([]byte("1"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("1"))
		}

		return nil
	})

	bc := Blockchain{tip, db}

	return &bc
}
