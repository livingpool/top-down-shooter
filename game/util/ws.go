package util

type ClientUpdate struct {
	PlayerId  string   `json:"player_id"`
	Type      string   `json:"type"`
	Keys      KeyPress `json:"keys"`
	Seq       int      `json:"seq"`
	TimeStamp int      `json:"timestamp"`
}

type KeyPress struct {
	W     bool `json:"w"`     // up
	S     bool `json:"s"`     // down
	A     bool `json:"a"`     // left
	D     bool `json:"d"`     // right
	Space bool `json:"space"` // shoot
}
