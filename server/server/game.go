package server

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/livingpool/top-down-shooter/game/util"
)

// saveClientUpdate takes a player's keystrokes and store them in the appropriate game world
func (gs *GameServer) saveClientUpdate(gameId uuid.UUID, update util.ClientUpdate) error {
	if g, exists := gs.games[gameId]; exists {
		id, err := uuid.Parse(update.PlayerId)
		if err != nil {
			return fmt.Errorf("error parsing player id: %s", update.PlayerId)
		}

		if p, exists := g.Players[id]; exists {
			p.ClientUpdates = append(p.ClientUpdates, update)
		} else {
			return fmt.Errorf("player id not found: %v", id)
		}
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
