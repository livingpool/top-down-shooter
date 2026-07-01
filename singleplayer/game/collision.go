package game

import (
	"log/slog"

	"github.com/livingpool/top-down-shooter/singleplayer/util"
)

// TODO: locality? only resolve collision for objects near player / zombie, rather than iterating thru all of them

func (g *Game) ResolveCollisions() {
	var maxPenVect util.Vector // resolve deepest collision (max penetration vector)
	var maxLen float64
	for _, obj := range g.Background.Objects {
		if vec, yes := g.Player.Object.Collide(*obj); yes {
			slog.Info("collision player <-> background obj", "position", g.Player.Object.Center)
			if l := vec.Length(); l > maxLen {
				maxPenVect = vec
				maxLen = l
			}
		}
	}

	if maxLen > 0 {
		g.Player.Object.Center.Sub(maxPenVect)
		g.Camera.Sub(maxPenVect)
		slog.Info("adjusted position of camera & player", "position", g.Player.Object.Center)
	}
}
