package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/livingpool/top-down-shooter/game/pkg/player"
	"github.com/livingpool/top-down-shooter/game/util"
)

type Game struct {
	DebugMode bool
	Players   []*player.Player
}

func (g *Game) Update() error {
	for _, p := range g.Players {
		p.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.Players {
		p.Draw(screen, g.DebugMode)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return util.ScreenWidth, util.ScreenHeight
}

// TODO: randomize players' spawn positions; at both initial spawn and reset
func (g *Game) Reset() {
	for i := range g.Players {
		g.Players[i] = player.NewPlayer(util.ScreenWidth, util.ScreenHeight)
	}
}
