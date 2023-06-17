package models

import (
	"encoding/json"
	"log"
)

// Server Message Types
const (
	Ping                  = 1000
	GameStartMessage      = 1001
	GameEndMessage        = 1002
	PositionUpdateMessage = 1003
	ScoreUpdateMessage    = 1004
	MatchFoundMessage     = 1005
)

// SocketMessage ..
type SocketMessage struct {
	Type    int             `json:"type"`
	Message json.RawMessage `json:"msg"`
}

// ScoreUpdate ..
type ScoreUpdate struct {
	Health    int `json:"health"`
	OppHealth int `json:"opp_health"`
}

// ParseSocketMessage ..
func ParseSocketMessage(msgBytes []byte) *SocketMessage {
	msg := SocketMessage{}
	err := json.Unmarshal(msgBytes, &msg)
	if err != nil {
		log.Println("ERROR", "Error in receiving message", err)
	}
	return &msg
}

// ToBytes returns the socket message in bytes
func (msg *SocketMessage) ToBytes() (returnMsg []byte) {
	returnMsg, err := json.Marshal(msg)
	if err != nil {
		returnMsg = []byte(err.Error())
	}
	return
}

func (msg *MatchFound) MatchFoundMessageToBytes() []byte {
	matchFoundBytes, err := json.Marshal(msg)
	if err != nil {
		matchFoundBytes = []byte(err.Error())
	}
	sm := SocketMessage{
		Type:    MatchFoundMessage,
		Message: matchFoundBytes,
	}
	return sm.ToBytes()
}

// GameStartMessageToBytes returns the game start message to socket message in bytes
func (msg *GameStart) GameStartMessageToBytes() []byte {
	gameStartBytes, err := json.Marshal(msg)
	if err != nil {
		gameStartBytes = []byte(err.Error())
	}
	sm := SocketMessage{
		Type:    GameStartMessage,
		Message: gameStartBytes,
	}
	return sm.ToBytes()
}

// GameEndMessageToSocket returns the game end message to socket message in bytes
func (msg *GameEnd) GameEndMessageToSocket() []byte {
	gameEndBytes, err := json.Marshal(msg)
	if err != nil {
		gameEndBytes = []byte(err.Error())
	}
	sm := SocketMessage{
		Type:    GameEndMessage,
		Message: gameEndBytes,
	}
	return sm.ToBytes()
}

// GetScoreUpdateSocketBytes ..
func GetScoreUpdateSocketBytes(health int, oppHealth int) []byte {
	su := ScoreUpdate{
		Health:    health,
		OppHealth: oppHealth,
	}
	suBytes, err := json.Marshal(su)
	if err != nil {
		suBytes = []byte(err.Error())
	}
	sm := SocketMessage{
		Type:    ScoreUpdateMessage,
		Message: suBytes,
	}
	return sm.ToBytes()
}
