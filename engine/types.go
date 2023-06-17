package engine

type WSMessage struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

type Login struct {
	ClientID string `json:"clientID"`
	Username string `json:"username"`
}
