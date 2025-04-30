package model

import (
	"github.com/google/uuid"
	"github.com/livingpool/top-down-shooter/game/util"
)

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

type Player struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	State    int     `json:"state"`
	Rotation float64 `json:"rotation"`
	Health   int     `json:"health"`
	Ammo     int     `json:"ammo"`
}

func NewPlayer(name string) *Player {
	id := uuid.New()
	return &Player{
		ID:       id.String(),
		Name:     name,
		X:        util.InitialPlayerX,
		Y:        util.InitialPlayerY,
		State:    util.HumanoidStateGun,
		Rotation: util.InitialPlayerRotation,
		Health:   util.InitialPlayerHealth,
		Ammo:     util.InitialPlayerAmmo,
	}
}

type Bullet struct {
	ID       string  `json:"id"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Rotation float64 `json:"rotation"`
}
