package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/DLzer/go-player-two/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Server Message Types
const (
	Ping                  = 1000
	GameStartMessage      = 1001
	GameEndMessage        = 1002
	PositionUpdateMessage = 1003
	ScoreUpdateMessage    = 1004
	wsServerEndpoint      = "ws://localhost:40000/engine"
	dataSendTickRate      = 100
)

func main() {
	dialer := websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	playerID := uuid.New().String()
	playerName := "dillon"
	connectionToServer := fmt.Sprintf("%s?id=%s_%s", wsServerEndpoint, playerName, playerID)

	conn, _, err := dialer.Dial(connectionToServer, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to socket")

	var msg models.SocketMessage
	for {
		_, m, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("WS read error", err)
		}
		if err := json.Unmarshal(m, &msg); err != nil {
			fmt.Println("WS unmarshal error", err)
		}

		switch msg.Type {
		case Ping:
			fmt.Println("received ping from server")
		case GameStartMessage:
			fmt.Println("received game start message")
			fmt.Printf("msg: %+v\n", msg.Message)
		case GameEndMessage:
			fmt.Println("received game end message")
			fmt.Printf("msg: %+v\n", msg.Message)
			conn.Close()
		case PositionUpdateMessage:
			fmt.Println("received position update message")
			fmt.Printf("msg: %+v\n", msg.Message)
		case ScoreUpdateMessage:
			fmt.Println("received score update message")
			fmt.Printf("msg: %+v\n", msg.Message)
		default:
			fmt.Println("received message we do not know")
			fmt.Printf("msg: %+v\n", msg.Message)
		}
	}
}
