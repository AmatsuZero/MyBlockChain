package main

import (
	"errors"
	"github.com/davecgh/go-spew/spew"
	"sync"
	"time"
)

type MyBlockchain struct {
	/* Chain is a series of validated Blocks*/
	Chain []Block
	mutex *sync.Mutex
}

func (bc *MyBlockchain) lastBlock() Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *MyBlockchain) isBlockValid(newBlock Block) bool {
	oldBlock := bc.lastBlock()
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if newBlock.calculateHash() != newBlock.Hash {
		return false
	}
	return true
}
func (bc *MyBlockchain) append(block Block) {
	bc.Chain = append(bc.Chain, block)
}

func (bc *MyBlockchain) replaceChain(newBlocks []Block) {
	bc.mutex.Lock()
	if len(newBlocks) > len(bc.Chain) {
		bc.Chain = newBlocks
	}
	bc.mutex.Unlock()
}

func (bc *MyBlockchain) generateBlock(BPM int) (Block, error) {
	var newBlock Block
	block := bc.lastBlock()

	t := time.Now()
	newBlock.Index = block.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = block.Hash
	newBlock.Hash = newBlock.calculateHash()

	return newBlock, nil
}

func (bc MyBlockchain) GenerateNewBlock(bpm int) (newBlock Block, err error) {
	newBlock, err = bc.generateBlock(bpm)
	if err != nil {
		return
	}
	if bc.isBlockValid(newBlock) {
		newBlockchain := append(bc.Chain, newBlock)
		bc.replaceChain(newBlockchain)
	} else {
		err = errors.New("not valid")
	}
	return
}

// Genesis create genesis block
func (bc *MyBlockchain) Genesis() {
	if bc.mutex != nil {
		return
	}
	bc.mutex = &sync.Mutex{}
	t := time.Now()
	genesisBlock := Block{0, t.String(), 0, "", ""}
	genesisBlock.Hash = genesisBlock.calculateHash()
	spew.Dump(genesisBlock)
	bc.append(genesisBlock)
}
