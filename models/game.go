package models

type GameStart struct {
	OpponentName string `json:"opponent"`
}

type GameEnd struct {
	Winner string `json:"winner"`
}
