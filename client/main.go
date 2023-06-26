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

// SocketMessage ..
type SocketMessage struct {
	Type    int             `json:"type"`
	Message json.RawMessage `json:"msg"`
}

// Position represents the individual players position in the game
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// PlayerState represents the individuals player state in the game
type PlayerState struct {
	Health     int      `json:"health"`
	Position   Position `json:"position"`
	IsOpponent bool     `json:"is_opp"`
}

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

	log.Println("Connected to socket")
	// liveChannel := make(chan int, 1)

	clientReader(conn)
	// clientWriter(conn, liveChannel)

	// close
}

func clientReader(conn *websocket.Conn) {
	var msg models.SocketMessage

BREAK:
	for {
		_, m, err := conn.ReadMessage()
		if err != nil {
			log.Println("WS read error", err)
		}
		if err := json.Unmarshal(m, &msg); err != nil {
			log.Println("WS unmarshal error", err)
		}

		var gameMessage interface{}

		fmt.Println("MESSAGE TYPE: ", msg.Type)

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
			break BREAK
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
			fmt.Printf("default msg: %+v\n", gameMessage)
			break BREAK
		}
	}
}

// func clientWriter(conn *websocket.Conn, liveChannel <-chan int) {
// 	for {
// 		running := <-liveChannel
// 		log.Println("Channel is open: ", running)
// 		if running == 0 {
// 			break
// 		}

// 		state := &PlayerState{
// 			Health: rand.Intn(100),
// 			Position: Position{
// 				X: rand.Intn(100),
// 				Y: rand.Intn(100),
// 			},
// 			IsOpponent: true,
// 		}

// 		j, err := json.Marshal(state)
// 		if err != nil {
// 			fmt.Println("Marshal Error:", err)
// 		}

// 		msg := &SocketMessage{
// 			Type:    PositionUpdateMessage,
// 			Message: j,
// 		}
// 		err = conn.WriteJSON(msg)
// 		if err != nil {
// 			fmt.Println("Write Error:", err)
// 		}

// 		time.Sleep(1 * time.Second)
// 	}
// 	return
// }
