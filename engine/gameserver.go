package engine

import (
	"log"
	"sync"
	"time"

	"github.com/DLzer/go-player-two/models"
)

type GameServer struct {
	clients     map[string]*Player
	games       map[string]*Game
	matchmaking chan *Player
	mu          sync.Mutex
}

var GS *GameServer

func init() {
	StartGameServer()
}

// StartGameServer initiates GameServer and starts the matchmaking listener
func StartGameServer() {
	GS = &GameServer{
		clients:     make(map[string]*Player),
		games:       make(map[string]*Game),
		matchmaking: make(chan *Player, 2),
	}

	go GS.MatchMaking()
}

// AddClient adds a new player to the client list
func (s *GameServer) AddClient(p *Player) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[p.Name] = p
}

// RemoveClient removes a player from the client list
func (s *GameServer) RemoveClient(playerName string, p *Player) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.clients, playerName)
}

// AddGame adds a new game to the games list
func (s *GameServer) AddGame(gameName string, g *Game) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.games[gameName] = g
}

// RemoveGame removes a game from the games list
func (s *GameServer) RemoveGame(gameName string, g *Game) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.games, gameName)
}

// GetStats returns a map of stats from the game server
func (s *GameServer) GetStats() map[string]int {
	return map[string]int{
		"games":   len(s.games),
		"players": len(s.clients),
		"status":  200,
	}
}

// MatchMaking will attempt to find a matching player connection
func (s *GameServer) MatchMaking() {
	log.Println("starting matchmaking")
	for {
		p1 := <-s.matchmaking
		p2 := <-s.matchmaking

		p1MatchFound := models.MatchFound{
			MatchFound: "match found",
		}

		p2MatchFound := models.MatchFound{
			MatchFound: "match found",
		}

		p1Err := p1.Send(p1MatchFound.MatchFoundMessageToBytes())
		p2Err := p2.Send(p2MatchFound.MatchFoundMessageToBytes())

		if p1Err != nil && p2Err != nil {
			continue
		} else if p1Err != nil {
			s.matchmaking <- p2
			continue
		} else if p2Err != nil {
			s.matchmaking <- p1
			continue
		} else {
			g := NewGame(p1, p2)
			s.AddGame("game_"+time.Now().String(), g)
			continue
		}
	}
}

// FindMatch will enter the player into the matchmaking pool
func (s *GameServer) FindMatch(p *Player) {
	s.AddClient(p)
	s.matchmaking <- p
}
