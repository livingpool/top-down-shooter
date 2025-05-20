package server

import (
	"context"
	"fmt"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	"github.com/livingpool/top-down-shooter/game/game"
	"github.com/livingpool/top-down-shooter/game/util"
)

// saveClientUpdate takes a player's keystrokes and store them in the appropriate game world
func (gs *GameServer) saveClientUpdate(game *game.Game, update util.ClientUpdate) error {
	id, err := uuid.Parse(update.PlayerId)
	if err != nil {
		return fmt.Errorf("error parsing player id: %s", update.PlayerId)
	}

	if p, exists := game.Players[id]; exists {
		p.ClientUpdates = append(p.ClientUpdates, update)
	} else {
		return fmt.Errorf("player id not found: %v", id)
	}

	return nil
}

// updatePhysics updates each game world at fixed deltas
func (gs *GameServer) updatePhysics() {
}

// TODO: rename this maybe
func (gs *GameServer) processInput(game *game.Game) error {
	for _, p := range game.Players {
		bullet := p.Update()
		if bullet != nil {
			game.Bullets[bullet.ID] = bullet
		}
	}

	return nil
}

func (gs *GameServer) sendServerUpdate(game *game.Game) error {
	for _, p := range game.Players {
		err := p.Conn.Write(context.TODO(), websocket.MessageText)
		if err != nil {
			return err
		}
	}
	return nil
}
