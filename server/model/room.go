package model

import (
	"sync"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type Room struct {
	ID          string
	Mutex       sync.Mutex
	Subscribers map[string]*Subscriber
	Broadcast   chan *ClientUpdate
}

func NewRoom() *Room {
	return &Room{
		ID:          uuid.NewString(),
		Mutex:       sync.Mutex{},
		Subscribers: make(map[string]*Subscriber),
		Broadcast:   make(chan *ClientUpdate),
	}
}

type Subscriber struct {
	Conn     *websocket.Conn
	Player   *Player
	Messages []ClientUpdate
}

type Action int

const (
	ActionRegisterClient = iota
	ActionUnregisterClient
)

type ClientRequest struct {
	Action     Action `json:"action"`
	RoomID     string `json:"room_id"`
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
}
