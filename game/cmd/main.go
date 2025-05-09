package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/livingpool/top-down-shooter/game/game"
	"github.com/livingpool/top-down-shooter/game/pkg/player"
)

func main() {
	g := &game.Game{
		DebugMode: true,
		Players:   []*player.Player{player.NewPlayer("todo")},
	}

	ebiten.SetWindowTitle("Tim's Top Down Shooter <3")

	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatalf("error running the game: %v", err)
	}
}
