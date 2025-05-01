package state

import (
	"github.com/livingpool/top-down-shooter/game/pkg/bullet"
	"github.com/livingpool/top-down-shooter/game/pkg/player"
	"github.com/livingpool/top-down-shooter/game/pkg/spawner"
)

type GameState struct {
	Players []*player.Player
	Bullets []*bullet.Bullet
	Zombies []*spawner.Zombie
}

type Msg struct {
}

type Ping struct {
}

type ServerUpdate struct {
}

type ClientUpdate struct {
	Type      string   `json:"type"`
	Keys      KeyPress `json:"keys"`
	TimeStamp int      `json:"timestamp"`
}

type KeyPress struct {
	Up    bool `json:"up"`
	Down  bool `json:"down"`
	Left  bool `json:"left"`
	Right bool `json:"right"`
	Shoot bool `json:"shoot"`
}
