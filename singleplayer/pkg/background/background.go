package background

import (
	"image"
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
	Objects []*util.GameObject
}

func NewBackground() *Background {
	trees := []*util.GameObject{
		util.NewGameObject(&util.Point{X: 300, Y: 100}, 0, assets.Tiles[182], util.CircleCollider),
		util.NewGameObject(&util.Point{X: 250, Y: 50}, 0, assets.Tiles[182], util.CircleCollider),
	}

	tilt := 3.0 * math.Pi / 180.0
	c, s := 64*math.Cos(tilt), 64*math.Sin(tilt) // for calculating the centers of the walls, c > s for tilt < 45 degrees

	var sx, sy float64 = 120, 200 // starting center of a corner sprite of the walls, all other walls are relative to it

	w, h := assets.Tiles[0].Bounds().Dx(), assets.Tiles[0].Bounds().Dy() // assume all images have same w and h

	walls := []*util.GameObject{
		// center block
		util.NewGameObject(&util.Point{X: sx, Y: sy}, tilt, assets.Tiles[108], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - s, Y: sy + c}, tilt, assets.Tiles[114], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + c, Y: sy + s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 2*c, Y: sy + 2*s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 3*c, Y: sy + 3*s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 4*c, Y: sy + 4*s}, tilt, assets.Tiles[139], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 4*c + s, Y: sy + 4*s - c}, tilt, assets.Tiles[140], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 4*c - s, Y: sy + 4*s + c}, tilt, assets.Tiles[135], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 5*c - s, Y: sy + 5*s + c}, tilt, assets.Tiles[113], util.RectCollider),

		// concrete floor at the top right corner
		util.NewGameObject(&util.Point{X: sx + 4*c + 2*s, Y: sy + 4*s - 2*c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 4*c + 2*s, Y: sy + 4*s - 2*c}, tilt, assets.Tiles[436], util.NoCollider), // glass door
		util.NewGameObject(&util.Point{X: sx + 5*c + 2*s, Y: sy + 5*s - 2*c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 6*c + 2*s, Y: sy + 6*s - 2*c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 7*c + 2*s, Y: sy + 7*s - 2*c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 5*c + s, Y: sy + 5*s - c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 6*c + s, Y: sy + 6*s - c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 7*c + s, Y: sy + 7*s - c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(
			&util.Point{X: sx + 7.5*c + s, Y: sy + 7.5*s - c}, // note that the position has to be adjusted as well
			tilt,
			assets.Tiles[10].SubImage(image.Rect(0, 0, w, h)).(*ebiten.Image), // show left half of the image
			util.NoCollider,
		),
		util.NewGameObject(&util.Point{X: sx + 5*c, Y: sy + 5*s}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 6*c, Y: sy + 6*s}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 7*c, Y: sy + 7*s}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(
			&util.Point{X: sx + 6*c - 0.5*s, Y: sy + 7*s + 0.5*c},
			tilt,
			assets.Tiles[10].SubImage(image.Rect(0, 0, w, h/2)).(*ebiten.Image), // show top half of the image
			util.NoCollider,
		),

		// wooden floor
		util.NewGameObject(&util.Point{X: sx - s + c, Y: sy + c + s}, tilt+math.Pi, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - s + 2*c, Y: sy + c + 2*s}, tilt+math.Pi, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - s + 3*c, Y: sy + c + 3*s}, tilt+math.Pi, assets.Tiles[41], util.NoCollider),
		// above 3 is to cover for a small gap between the kitchen area and the wooden floor
		util.NewGameObject(&util.Point{X: sx - 2*s, Y: sy + 2*c}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 2*s, Y: sy + 2*c}, tilt, assets.Tiles[436], util.NoCollider), // glass door
		util.NewGameObject(&util.Point{X: sx - 3*s, Y: sy + 3*c}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 3*s, Y: sy + 3*c}, tilt, assets.Tiles[436], util.NoCollider), // glass door
		util.NewGameObject(&util.Point{X: sx - 3*s + c, Y: sy + 3*c + s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 3*s + 2*c, Y: sy + 3*c + 2*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 3*s + 3*c, Y: sy + 3*c + 3*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 3*s + 4*c, Y: sy + 3*c + 4*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 3*s + 5*c, Y: sy + 3*c + 5*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 3*s + 6*c, Y: sy + 3*c + 6*s}, tilt, assets.Tiles[43], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 3*s + 7*c, Y: sy + 3*c + 7*s}, tilt, assets.Tiles[45], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 3*s + 8*c, Y: sy + 3*c + 8*s}, tilt, assets.Tiles[45], util.NoCollider),

		util.NewGameObject(&util.Point{X: sx - 2*s + c, Y: sy + 2*c + s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 2*s + 2*c, Y: sy + 2*c + 2*s}, tilt, assets.Tiles[43], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 2*s + 3*c, Y: sy + 2*c + 3*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 2*s + 4*c, Y: sy + 2*c + 4*s}, tilt, assets.Tiles[45], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 2*s + 5*c, Y: sy + 2*c + 5*s}, tilt, assets.Tiles[45], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 2*s + 6*c, Y: sy + 2*c + 6*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 2*s + 7*c, Y: sy + 2*c + 7*s}, tilt, assets.Tiles[43], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 2*s + 8*c, Y: sy + 2*c + 8*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 2*s + 8*c, Y: sy + 2*c + 8*s}, tilt, assets.Tiles[437], util.NoCollider), // glass door

		util.NewGameObject(
			&util.Point{X: sx + 6*c - 1.5*s, Y: sy + 5*s + 1.5*c},
			tilt,
			assets.Tiles[41].SubImage(image.Rect(0, 0, w, h/2)).(*ebiten.Image), // show top half of the image
			util.NoCollider,
		),

		util.NewGameObject(&util.Point{X: sx - 4*s + c, Y: sy + 4*c + s}, tilt, assets.Tiles[43], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 4*s + 2*c, Y: sy + 4*c + 2*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 4*s + 3*c, Y: sy + 4*c + 3*s}, tilt, assets.Tiles[45], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 4*s + 4*c, Y: sy + 4*c + 4*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 4*s + 5*c, Y: sy + 4*c + 5*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 4*s + 6*c, Y: sy + 4*c + 6*s}, tilt, assets.Tiles[45], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 4*s + 7*c, Y: sy + 4*c + 7*s}, tilt, assets.Tiles[43], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 4*s + 8*c, Y: sy + 4*c + 8*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 4*s + 8*c, Y: sy + 4*c + 8*s}, tilt, assets.Tiles[437], util.NoCollider), // glass door

		// kitchen area
		util.NewGameObject(&util.Point{X: sx - s + c, Y: sy + c + s}, tilt+math.Pi, assets.Tiles[323], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - s + 2*c, Y: sy + c + 2*s}, tilt+math.Pi, assets.Tiles[322], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - s + 3*c, Y: sy + c + 3*s}, tilt+math.Pi, assets.Tiles[320], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - s + 3*c, Y: sy + c + 3*s}, tilt+math.Pi, assets.Tiles[267], util.NoCollider),

		// right middle corner
		util.NewGameObject(&util.Point{X: sx + 7*c - s, Y: sy + 7*s + c}, tilt, assets.Tiles[141], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 8*c - s, Y: sy + 8*s + c}, tilt, assets.Tiles[136], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 8*c, Y: sy + 8*s}, tilt, assets.Tiles[140], util.RectCollider),

		// upper right corner
		util.NewGameObject(&util.Point{X: sx + 4*c + 3*s, Y: sy + 4*s - 3*c}, tilt, assets.Tiles[141], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 5*c + 3*s, Y: sy + 5*s - 3*c}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 6*c + 3*s, Y: sy + 6*s - 3*c}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 7*c + 3*s, Y: sy + 7*s - 3*c}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 8*c + 3*s, Y: sy + 8*s - 3*c}, tilt, assets.Tiles[109], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 8*c + 2*s, Y: sy + 8*s - 2*c}, tilt, assets.Tiles[114], util.RectCollider),

		// wooden doors
		util.NewGameObject(&util.Point{X: sx + 6*c - s, Y: sy + 6*s + c}, tilt, assets.Tiles[467], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 8*c + s, Y: sy + 8*s - c}, tilt, assets.Tiles[440], util.NoCollider),

		// pillars
		{
			Center:   &util.Point{X: sx + 4*c - 3*s, Y: sy + 4*s + 3*c},
			Rotation: tilt,
			Sprite:   assets.Tiles[170],
			Collider: util.NewRect(&util.Point{X: sx + 4*c - 3*s, Y: sy + 4*s + 3*c}, 48, 48, tilt),
		},
		util.NewGameObject(&util.Point{X: sx + 8*c - 3*s, Y: sy + 8*s + 3*c}, tilt, assets.Tiles[195], util.RectCollider),

		// bottom block
		util.NewGameObject(&util.Point{X: sx - 4*s, Y: sy + 4*c}, tilt, assets.Tiles[140], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 5*s, Y: sy + 5*c}, tilt, assets.Tiles[135], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 5*s + c, Y: sy + 5*c + s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 5*s + 2*c, Y: sy + 5*c + 2*s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 5*s + 3*c, Y: sy + 5*c + 3*s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 5*s + 4*c, Y: sy + 5*c + 4*s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 5*s + 5*c, Y: sy + 5*c + 5*s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 5*s + 6*c, Y: sy + 5*c + 6*s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 5*s + 7*c, Y: sy + 5*c + 7*s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 5*s + 8*c, Y: sy + 5*c + 8*s}, tilt, assets.Tiles[113], util.RectCollider),

		// other objects without colliders
		util.NewGameObject(&util.Point{X: sx + 7*c + 4*s, Y: sy + 5*s - 2*c}, tilt+10.0*math.Pi/180.0, assets.Tiles[129], util.NoCollider),

		// other objects with colliders
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

	for _, obj := range b.Objects {
		op := obj.CenterAndRotateImage()
		screen.DrawImage(obj.Sprite, op)

		// if debugMode {
		// 	switch c := obj.Collider.(type) {
		// 	case util.Rect:
		// 		obj.DrawDebugRect(screen, float32(c.DimX), float32(c.DimY), "")
		// 	case util.Circle:
		// 		obj.DrawDebugCircle(screen, float32(c.Radius), "tree")
		// 	}
		// }
	}
}
