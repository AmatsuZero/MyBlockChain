package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"time"
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

type Message struct {
	BPM int
}

// Blockchain is a series of validated Blocks
type Blockchain struct {
	blockChain []Block
}

var blockChain Blockchain

func (bc *Blockchain) lastBlock() Block {
	return bc.blockChain[len(bc.blockChain)-1]
}

func (bc *Blockchain) isBlockValid(newBlock Block) bool {
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
func (bc *Blockchain) append(block Block) {
	bc.blockChain = append(bc.blockChain, block)
}

func (bc *Blockchain) replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(bc.blockChain) {
		bc.blockChain = newBlocks
	}
}

func (bc *Blockchain) generateBlock(BPM int) (Block, error) {
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

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(blockChain.blockChain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Message
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()
	newBlock, err := blockChain.generateBlock(m.BPM)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	if blockChain.isBlockValid(newBlock) {
		newBlockchain := append(blockChain.blockChain, newBlock)
		blockChain.replaceChain(newBlockchain)
		spew.Dump(blockChain)
	}
	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func run() error {
	muxRouter := makeMuxRouter()
	httpAddr := os.Getenv("PORT")
	log.Println("Listening on ", os.Getenv("PORT"))
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        muxRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		t := time.Now()
		genesisBlock := Block{0, t.String(), 0, "", ""}
		genesisBlock.Hash = genesisBlock.calculateHash()
		spew.Dump(genesisBlock)
		blockChain.append(genesisBlock)
	}()
	log.Fatal(run())
}
