package background

import (
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/livingpool/top-down-shooter/game/assets"
	"github.com/livingpool/top-down-shooter/singleplayer/util"
)

// So the idea is:
// Initially, draw the background wrt the screen.
// Additional tiles will be drawn outside the viewport.
// Player will be at the center of the screen.
//
// When the player moves, track the [changes in coord].
// Draw every image of the background wrt this change.
// Keep the top-left-most tile at max (0, 0).
// If this is reached, dont change background further. This signals the boundary.
//
// On second thought, it would be best to keep the camera with the player!
//
// Finally, draw the player.

type Background struct {
	Objects []*util.GameObject // their relative positions are fixed
}

func NewBackground() *Background {
	// TODO: next step: draw the map!
	trees := []*util.GameObject{
		util.NewGameObject(&util.Point{X: 100, Y: 100}, 0, assets.SmallBrownTree, util.CircleCollider),
	}

	return &Background{
		Objects: slices.Concat(trees),
	}
}

// b.Update updates every background objects' centers and their colliders' centers
func (b *Background) Update() {
}

// b.Draw the background w.r.t the camera (player)
func (b *Background) Draw(screen *ebiten.Image, offsetX, offsetY float64, debugMode bool) {
	repeat := max(util.ScreenWidth/assets.Tile1.Bounds().Dx(), util.ScreenHeight/assets.Bullet.Bounds().Dy()) + 2

	w, h := assets.Tile1.Bounds().Dx(), assets.Tile1.Bounds().Dy()

	for i := -1; i <= repeat; i++ {
		for j := -1; j <= repeat; j++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(2*w*i), float64(2*h*j))
			op.GeoM.Translate(offsetX, offsetY)

			screen.DrawImage(assets.Tile1, op)

			op.GeoM.Translate(float64(w), 0)
			screen.DrawImage(assets.Tile2, op)

			op.GeoM.Translate(0, float64(h))
			screen.DrawImage(assets.Tile3, op)

			op.GeoM.Translate(-float64(w), 0)
			screen.DrawImage(assets.Tile4, op)
		}
	}

	// objects that are subject to collision
	for _, obj := range b.Objects {
		op := obj.CenterAndRotateImage()
		screen.DrawImage(obj.Sprite, op)

		if debugMode {
			obj.DrawDebugCircle(screen, 32, "tree")
		}
	}
}
