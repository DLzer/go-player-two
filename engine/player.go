package engine

import (
	"log"
	"sync"
	"time"

	"github.com/DLzer/go-player-two/models"
	"github.com/gorilla/websocket"
)

// Player represents an individual connected to the server
type Player struct {
	conn      *websocket.Conn
	Message   chan []byte
	Exit      chan int
	sendMutex sync.Mutex
	Name      string `json:"name"`
	Ready     bool
	State     models.PlayerState
	Position  models.Position
}

// SpawnNewPlayer returns a new player instance
func SpawnNewPlayer(playerName string, conn *websocket.Conn) *Player {
	return &Player{
		conn:    conn,
		Message: make(chan []byte, 5),
		Exit:    make(chan int),
		Name:    playerName,
	}
}

// Receive receives messages and or an exit signal from the client
func (p *Player) Receive() {
	for {
		_, message, err := p.conn.ReadMessage()
		p.Message <- message
		if err != nil {
			log.Println(p.Name, " has quit.")
			p.conn.Close()
			// notify that the player has quit
			close(p.Exit)
			break
		}
	}
}

// Send will send a message to the player
func (p *Player) Send(msg []byte) error {
	p.sendMutex.Lock()
	defer p.sendMutex.Unlock()

	return p.conn.WriteMessage(1, msg)
}

// Close will exit the player session and close the connection
func (p *Player) Close() {
	time.Sleep(1 * time.Second)
	p.conn.Close()
}
