package engine

import (
	"encoding/json"
	"log"

	"github.com/DLzer/go-player-two/models"
)

const (
	StartingHealth = 100
)

type Game struct {
	ID        string
	PlayerOne *Player
	PlayerTwo *Player
}

func NewGame(gameID string, p1 *Player, p2 *Player) *Game {
	log.Println("starting game instance: ", gameID)
	p1.State.Health = StartingHealth
	p2.State.Health = StartingHealth
	game := Game{
		ID:        gameID,
		PlayerOne: p1,
		PlayerTwo: p2,
	}
	// Game is starting, check player connections
	go game.RouteMessage()
	// go game.RouteMessage()

	// Sending game start signal
	initializePlayerOne := models.GameStart{
		OpponentName: p2.Name,
	}
	if err := p1.Send(initializePlayerOne.GameStartMessageToBytes()); err != nil {
		log.Fatal(err)
	}

	// Sending game start signal
	initializePlayerTwo := models.GameStart{
		OpponentName: p1.Name,
	}
	if err := p2.Send(initializePlayerTwo.GameStartMessageToBytes()); err != nil {
		log.Fatal(err)
	}

	return &game
}

// RouteMessage ..
func (g *Game) RouteMessage() {
	log.Println("Starting RouteMessage ...")
BREAK:
	for {
		select {
		case msgBytes := <-g.PlayerOne.Message:
			msg := models.ParseSocketMessage(msgBytes)
			g.handleGameMove(msg, g.PlayerOne, g.PlayerTwo)

		case msgBytes := <-g.PlayerTwo.Message:
			msg := models.ParseSocketMessage(msgBytes)
			g.handleGameMove(msg, g.PlayerTwo, g.PlayerOne)

		// Handle exit messages
		case <-g.PlayerOne.Exit:
			iP1 := models.GameEnd{
				Winner: g.PlayerTwo.Name,
			}
			log.Println(g.PlayerOne.Name, " has exited the game")
			if err := g.PlayerTwo.Send(iP1.GameEndMessageToSocket()); err != nil {
				log.Fatal(err)
			}
			g.SafeShutdown()
			break BREAK

		case <-g.PlayerTwo.Exit:
			iP2 := models.GameEnd{
				Winner: g.PlayerOne.Name,
			}
			log.Println(g.PlayerTwo.Name, " has exited the game")
			if err := g.PlayerOne.Send(iP2.GameEndMessageToSocket()); err != nil {
				log.Fatal(err)
			}
			g.SafeShutdown()
			break BREAK
		}
	}
	log.Println("Closing RouteMessage game for", g.PlayerOne.Name, " and ", g.PlayerTwo.Name)
}

// SafeShutdown safely handles removing the clients and game from the mapping pool
func (g *Game) SafeShutdown() {
	go g.PlayerOne.Close()
	go g.PlayerTwo.Close()
	GS.RemoveClient(g.PlayerTwo.Name, g.PlayerTwo)
	GS.RemoveClient(g.PlayerOne.Name, g.PlayerOne)
	GS.RemoveGame(g.ID, g)
}

// handleGameMove ..
func (g *Game) handleGameMove(msg *models.SocketMessage, player *Player, opponent *Player) {
	switch msg.Type {
	case models.Ping:
		// Ping the message back to client
		if err := player.Send(msg.ToBytes()); err != nil {
			log.Fatal(err)
		}

	case models.PositionUpdateMessage:
		handlePositionUpdate(msg, player, opponent)
	}
}

func handlePositionUpdate(msg *models.SocketMessage, player *Player, opponent *Player) {
	var state models.PlayerState
	err := json.Unmarshal(msg.Message, &state)
	if err != nil {
		log.Println("ERROR", "Invalid message", msg.Message, err)
	}
	if !player.Ready {
		// The player is ready for the match
		log.Println(player.Name, "is ready for the match")
		player.Ready = true
	}

	// Update the player position value in server
	player.State.Position.X = state.Position.Y

	// Marshal the json and send to player
	j, _ := json.Marshal(state)
	msg.Message = j
	if err := player.Send(msg.ToBytes()); err != nil {
		log.Fatal(err)
	}

	// Marshal the json and send to opponent
	state.IsOpponent = true
	j, _ = json.Marshal(state)
	msg.Message = j
	if err := opponent.Send(msg.ToBytes()); err != nil {
		log.Fatal(err)
	}
}
