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
	MatchFoundMessage     = 1005
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

		var gameMessage interface{}

		switch msg.Type {
		case Ping:
			fmt.Println("received ping from server")
		case MatchFoundMessage:
			fmt.Println("received match found message")
			if err := json.Unmarshal(m, &gameMessage); err != nil {
				fmt.Println("WS unmarshal error", err)
			}
			fmt.Printf("msg: %+v\n", gameMessage)
		case GameStartMessage:
			fmt.Println("received game start message")
			if err := json.Unmarshal(m, &gameMessage); err != nil {
				fmt.Println("WS unmarshal error", err)
			}
			fmt.Printf("msg: %+v\n", gameMessage)
		case GameEndMessage:
			fmt.Println("received game end message")
			if err := json.Unmarshal(m, &gameMessage); err != nil {
				fmt.Println("WS unmarshal error", err)
			}
			fmt.Printf("msg: %+v\n", gameMessage)
			conn.Close()
			break
		case PositionUpdateMessage:
			if err := json.Unmarshal(m, &gameMessage); err != nil {
				fmt.Println("WS unmarshal error", err)
			}
			fmt.Printf("msg: %+v\n", gameMessage)
		case ScoreUpdateMessage:
			if err := json.Unmarshal(m, &gameMessage); err != nil {
				fmt.Println("WS unmarshal error", err)
			}
			fmt.Printf("msg: %+v\n", gameMessage)
		default:
			if err := json.Unmarshal(m, &gameMessage); err != nil {
				fmt.Println("WS unmarshal error", err)
			}
			fmt.Printf("msg: %+v\n", gameMessage)
		}
	}
}
