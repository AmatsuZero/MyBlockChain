package main

import (
	"bufio"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

type BlockChainServer struct {
	BlockChain *MyBlockchain
	// bcServer handles incoming concurrent Blocks
	bcServer chan []Block
}

func NewBlockChainServer() *BlockChainServer {
	var blockChain MyBlockchain
	blockChain.Genesis()
	return &BlockChainServer{
		BlockChain: &blockChain,
		bcServer:   make(chan []Block),
	}
}

func (s *BlockChainServer) StartAndListen() {
	tcpPort := os.Getenv("PORT")
	server, err := net.Listen("tcp", ":"+tcpPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("TCP  Server Listening on port :", tcpPort)
	defer server.Close()
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go s.handleConn(conn)
	}
}

func (s *BlockChainServer) handleConn(conn net.Conn) {
	defer conn.Close()

	io.WriteString(conn, "Enter a new BPM:")

	scanner := bufio.NewScanner(conn)

	// take in BPM from stdin and add it to blockchain after conducting necessary validation
	go func() {
		for scanner.Scan() {
			bpm, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Printf("%v not a number: %v", scanner.Text(), err)
				continue
			}
			_, err = s.BlockChain.GenerateNewBlock(bpm)
			if err != nil {
				log.Println(err)
				continue
			}
			s.bcServer <- s.BlockChain.Chain
			io.WriteString(conn, "\nEnter a new BPM:")
		}
	}()

	// simulate receiving broadcast
	go func() {
		for {
			time.Sleep(30 * time.Second)
			s.BlockChain.mutex.Lock()
			output, err := json.Marshal(s.BlockChain.Chain)
			if err != nil {
				log.Fatal(err)
			}
			s.BlockChain.mutex.Unlock()
			io.WriteString(conn, string(output))
		}
	}()

	for _ = range s.bcServer {
		spew.Dump(s.BlockChain.Chain)
	}
}
