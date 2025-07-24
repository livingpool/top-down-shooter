package background

import (
	"log"
	"math"
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
	trees := []*util.GameObject{
		util.NewGameObject(&util.Point{X: 300, Y: 100}, 0, assets.Tiles[182], util.CircleCollider),
		util.NewGameObject(&util.Point{X: 250, Y: 50}, 0, assets.Tiles[182], util.CircleCollider),
	}

	tilt := 3.0 * math.Pi / 180.0
	c, s := 64*math.Cos(tilt), 64*math.Sin(tilt) // for calculating the centers of the walls, c > s for tilt < 45 degrees

	var sx, sy float64 = 120, 200 // starting center of a corner sprite of the walls, all other walls are relative to it

	walls := []*util.GameObject{
		util.NewGameObject(&util.Point{X: sx, Y: sy}, tilt, assets.Tiles[108], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - s, Y: sy + c}, tilt, assets.Tiles[114], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + c, Y: sy + s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 2*c, Y: sy + 2*s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 3*c, Y: sy + 3*s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 4*c, Y: sy + 4*s}, tilt, assets.Tiles[139], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 4*c + s, Y: sy + 4*s - c}, tilt, assets.Tiles[140], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 4*c - s, Y: sy + 4*s + c}, tilt, assets.Tiles[135], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 5*c - s, Y: sy + 5*s + c}, tilt, assets.Tiles[113], util.RectCollider),
	}

	return &Background{
		Objects: slices.Concat(trees, walls),
	}
}

// b.Update updates every background objects' centers and their colliders' centers
func (b *Background) Update() {
}

// b.Draw draws the background w.r.t the camera (player)
func (b *Background) Draw(screen *ebiten.Image, offsetX, offsetY float64, debugMode bool) {
	w, h := assets.Tiles[0].Bounds().Dx(), assets.Tiles[0].Bounds().Dy() // 64, 64
	repeat := max(util.ScreenWidth/w, util.ScreenHeight/h) + 2

	// grass floor
	for i := -1; i <= repeat; i++ {
		for j := -1; j <= repeat; j++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(2*w*i), float64(2*h*j))
			op.GeoM.Translate(offsetX, offsetY)

			screen.DrawImage(assets.Tiles[0], op)

			op.GeoM.Translate(float64(w), 0)
			screen.DrawImage(assets.Tiles[1], op)

			op.GeoM.Translate(0, float64(h))
			screen.DrawImage(assets.Tiles[2], op)

			op.GeoM.Translate(-float64(w), 0)
			screen.DrawImage(assets.Tiles[3], op)
		}
	}

	// objects that are subject to collision
	for _, obj := range b.Objects {
		op := obj.CenterAndRotateImage()
		screen.DrawImage(obj.Sprite, op)

		if debugMode {
			switch obj.GetColliderType() {
			case util.RectCollider:
				obj.DrawDebugRect(screen, float32(w), float32(h), "")
			case util.CircleCollider:
				obj.DrawDebugCircle(screen, 32, "tree")
			default:
				log.Fatal("wtf bro")
			}
		}
	}
}
