package blockchain

import (
	"sync"

	"github.com/BlockChain/db"
	"github.com/BlockChain/utils"
)

//If any variable or function name is starting with capital letter
//it will be exported

//func (reference or reciver function) functionName(parameter) returnType {}

type blockChain struct {
	NewstHash string `json:"newestHash"`
	Height    int    `json:"height"`
}

var b *blockChain
var once sync.Once

func (b *blockChain) persist() {
	db.SaveBlockChain(utils.ToBytes(b))
}

func (b *blockChain) AddBlock(data string) {
	block := createBlock(data, b.NewstHash, b.Height+1)
	b.NewstHash = block.Hash
	b.Height = block.Height
	b.persist()
}

//Getting a blockchain
//If nothing is in the block chain, create a Genesis Block as the first block
func BlockChain() *blockChain {
	if b == nil {
		//Just do one time even though we are running thousands of goroutine
		//It will run only initialization time and this is called a Singleton Pattern
		once.Do(func() {
			b = &blockChain{"", 0}
			b.AddBlock("Genesis Block")
		})
	}
	return b
}
