package game

import "log/slog"

func (g *Game) ResolveCollisions() {
	for _, obj := range g.Background.Objects {
		if vec, yes := g.Player.Object.Collide(*obj); yes {
			slog.Debug("collision player <-> background obj", "position", g.Player.Object.Center)
			g.Player.Object.Center.Sub(vec)
		}
	}
}
