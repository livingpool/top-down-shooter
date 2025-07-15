package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/livingpool/top-down-shooter/singleplayer/game"
)

func main() {
	g := game.NewGame(true)

	ebiten.SetWindowTitle("Tim's Top Down Shooter <3")

	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatalf("error running the game: %v", err)
	}
}
