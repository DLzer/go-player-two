package models

type GameStart struct {
	OpponentName string `json:"opponent"`
}

type GameEnd struct {
	Winner string `json:"winner"`
}

type MatchFound struct {
	MatchFound string `json:"match_found"`
}
