package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DLzer/go-player-two/engine"
	"github.com/gorilla/websocket"
)

func main() {

	var addr = flag.String("addr", "40000", "http service address")

	http.HandleFunc("/engine", engineHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "running")
	})

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", *addr),
		ReadHeaderTimeout: 3 * time.Second,
	}

	fmt.Printf("Server starting on port %s\n", *addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	EnableCompression: true,
}

func engineHandler(w http.ResponseWriter, r *http.Request) {
	playerName := r.URL.Query()["id"][0]
	fmt.Printf("player joined the server: %s\n", playerName)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Spawn a new player
	p := engine.SpawnNewPlayer(playerName, conn)
	// Start listening to messages from player
	go p.Receive()
	// Find a match for the player
	engine.GS.FindMatch(p)
}
