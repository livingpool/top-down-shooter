package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/livingpool/top-down-shooter/game/pkg/player"
	"github.com/livingpool/top-down-shooter/game/util"
)

type Game struct {
	debugMode bool
	player    *player.Player
}

func (g *Game) Update() error {
	g.player.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.player.Draw(screen, g.debugMode)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return util.ScreenWidth, util.ScreenHeight
}

func (g *Game) Reset() {
	g.player = player.NewPlayer(util.ScreenWidth, util.ScreenHeight)
}

func main() {
	g := &Game{
		debugMode: true,
		player:    player.NewPlayer(util.ScreenWidth, util.ScreenHeight),
	}

	ebiten.SetWindowTitle("Tim's Top Down Shooter <3")

	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatalf("error running the game: %v", err)
	}
}
