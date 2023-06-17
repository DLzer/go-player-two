package models

// Position represents the individual players position in the game
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// PlayerState represents the individuals player state in the game
type PlayerState struct {
	Health     float64  `json:"health"`
	Position   Position `json:"position"`
	IsOpponent bool     `json:"is_opp"`
}
