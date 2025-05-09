package game

import "log"

// TODO: this obviously needs some fix
func (g *Game) checkCollisions() {
	// player collides with bullet
	for _, player := range g.Players {
		for _, b := range g.Bullets {
			if player.Collider().Intersects(b.Collider()) {
				log.Println("player collided with bullet!")
				// g.Bullets = slices.Delete(g.Bullets, j, j+1)
				// player.Health--
			}
		}
	}

	// player collides with player

	// player collides with zombie

	// zombie collides with zombie

	// zombie collides with bullet
}
