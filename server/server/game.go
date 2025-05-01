package server

import (
	"fmt"

	"github.com/livingpool/top-down-shooter/server/model"
)

// saveClientUpdate takes a player's keystrokes and store them in the appropriate game world
func (gs *GameServer) saveClientUpdate(gameId string, update model.ClientUpdate) error {
	if q, exists := gs.inputs[gameId]; exists {
		q = append(q, update)
	} else {
		return fmt.Errorf("game id not found: %v", gameId)
	}

	return nil
}

// updatePhysics updates each game world at fixed deltas
func (gs *GameServer) updatePhysics() {
}

func (gs *GameServer) processInput() {
}
