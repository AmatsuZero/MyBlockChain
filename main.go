package main

import (
	"github.com/joho/godotenv"
	"log"
)

type Message struct {
	BPM int
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	server := NewBlockChainServer()
	server.StartAndListen()
}
