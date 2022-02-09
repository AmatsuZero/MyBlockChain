package main

import (
	"crypto/sha256"
	"encoding/hex"
)

/*Block represents each 'item' in the blockchain*/
type Block struct {
	/*Index is the position of the data record in the blockchain*/
	Index int
	/*Timestamp is automatically determined and is the time the data is written*/
	Timestamp string
	/*BPM or beats per minute, is your pulse rate*/
	BPM int
	/*Hash is a SHA256 identifier representing this data record*/
	Hash string
	/*PrevHash is the SHA256 identifier of the previous record in the chain*/
	PrevHash string
}

func (b *Block) calculateHash() string {
	record := string(rune(b.Index)) + b.Timestamp + string(rune(b.BPM)) + b.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}
