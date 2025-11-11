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

// some constants for reuse
var (
	rotate180 = 180 * math.Pi / 180.0
	rotate90  = 90 * math.Pi / 180.0
	rotate45  = 45 * math.Pi / 180.0
	tilt      = 3.0 * math.Pi / 180.0
	c, s      = 64 * math.Cos(tilt), 64 * math.Sin(tilt) // for calculating the centers of the walls, c > s for tilt < 45 degrees

	sx, sy float64 = 120, 200 // starting center of a corner sprite of the walls, all other walls are relative to it

	w, h = assets.Tiles[0].Bounds().Dx(), assets.Tiles[0].Bounds().Dy() // assume all images have same w and h
)

// chopped images; note that their positions must be adjusted accordingly
var (
	topHalfConcreteFloor  = assets.Tiles[10].SubImage(image.Rect(0, 0, w, h/2)).(*ebiten.Image)
	leftHalfConcreteFloor = assets.Tiles[10].SubImage(image.Rect(0, 0, w, h)).(*ebiten.Image)
	topHalfWoldenFloor    = assets.Tiles[41].SubImage(image.Rect(0, 0, w, h/2)).(*ebiten.Image)
)

func NewBackground() *Background {
	trees := []*util.GameObject{
		util.NewGameObject(&util.Point{X: 300, Y: 100}, 0, assets.Tiles[182], util.CircleCollider),
		util.NewGameObject(&util.Point{X: 250, Y: 50}, 0, assets.Tiles[182], util.CircleCollider),
	}

	walls := []*util.GameObject{
		// center block
		util.NewGameObject(&util.Point{X: sx, Y: sy}, tilt, assets.Tiles[108], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - s, Y: sy + c}, tilt, assets.Tiles[114], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + c, Y: sy + s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 2*c, Y: sy + 2*s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 3*c, Y: sy + 3*s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 4*c, Y: sy + 4*s}, tilt, assets.Tiles[139], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 4*c + s, Y: sy + 4*s - c}, tilt, assets.Tiles[137], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 4*c + 2*s, Y: sy + 4*s - 2*c}, tilt, assets.Tiles[140], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 4*c - s, Y: sy + 4*s + c}, tilt, assets.Tiles[135], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 5*c - s, Y: sy + 5*s + c}, tilt, assets.Tiles[113], util.RectCollider),

		// concrete floor at the top right corner
		util.NewGameObject(&util.Point{X: sx + 4*c + 4*s, Y: sy + 4*s - 4*c}, tilt, assets.Tiles[10], util.NoCollider),  // start of 1st row
		util.NewGameObject(&util.Point{X: sx + 4*c + 4*s, Y: sy + 4*s - 4*c}, tilt, assets.Tiles[436], util.NoCollider), // glass door
		util.NewGameObject(&util.Point{X: sx + 5*c + 4*s, Y: sy + 5*s - 4*c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 6*c + 4*s, Y: sy + 6*s - 4*c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 7*c + 4*s, Y: sy + 7*s - 4*c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 8*c + 4*s, Y: sy + 8*s - 4*c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 9*c + 4*s, Y: sy + 9*s - 4*c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 4*c + 3*s, Y: sy + 4*s - 3*c}, tilt, assets.Tiles[10], util.NoCollider),  // start of 2nd row
		util.NewGameObject(&util.Point{X: sx + 4*c + 3*s, Y: sy + 4*s - 3*c}, tilt, assets.Tiles[436], util.NoCollider), // glass door
		util.NewGameObject(&util.Point{X: sx + 5*c + 3*s, Y: sy + 5*s - 3*c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 6*c + 3*s, Y: sy + 6*s - 3*c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 7*c + 3*s, Y: sy + 7*s - 3*c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 8*c + 3*s, Y: sy + 8*s - 3*c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 9*c + 3*s, Y: sy + 9*s - 3*c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 5*c + 2*s, Y: sy + 5*s - 2*c}, tilt, assets.Tiles[10], util.NoCollider), // start of 3rd row
		util.NewGameObject(&util.Point{X: sx + 6*c + 2*s, Y: sy + 6*s - 2*c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 7*c + 2*s, Y: sy + 7*s - 2*c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 8*c + 2*s, Y: sy + 8*s - 2*c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 9*c + 2*s, Y: sy + 9*s - 2*c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 9.5*c + 2*s, Y: sy + 9.5*s - 2*c}, tilt, leftHalfConcreteFloor, util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 5*c + s, Y: sy + 5*s - c}, tilt, assets.Tiles[10], util.NoCollider), // start of 4th row
		util.NewGameObject(&util.Point{X: sx + 6*c + s, Y: sy + 6*s - c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 7*c + s, Y: sy + 7*s - c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 8*c + s, Y: sy + 8*s - c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 9*c + s, Y: sy + 9*s - c}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 5*c, Y: sy + 5*s}, tilt, assets.Tiles[10], util.NoCollider), // start of 5th row
		util.NewGameObject(&util.Point{X: sx + 6*c, Y: sy + 6*s}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 7*c, Y: sy + 7*s}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 8*c, Y: sy + 8*s}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 9*c, Y: sy + 9*s}, tilt, assets.Tiles[10], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 6*c - 0.5*s, Y: sy + 7*s + 0.5*c}, tilt, topHalfConcreteFloor, util.NoCollider), // start of 6th row
		util.NewGameObject(&util.Point{X: sx + 7*c - 0.5*s, Y: sy + 8*s + 0.5*c}, tilt, topHalfConcreteFloor, util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 8*c - 0.5*s, Y: sy + 9*s + 0.5*c}, tilt, topHalfConcreteFloor, util.NoCollider),

		// wooden floor
		util.NewGameObject(&util.Point{X: sx - s + c, Y: sy + c + s}, tilt+math.Pi, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - s + 2*c, Y: sy + c + 2*s}, tilt+math.Pi, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - s + 3*c, Y: sy + c + 3*s}, tilt+math.Pi, assets.Tiles[41], util.NoCollider),
		// above 3 is to cover for a small gap between the kitchen area and the wooden floor
		util.NewGameObject(&util.Point{X: sx - 2*s, Y: sy + 2*c}, tilt, assets.Tiles[41], util.NoCollider),  // start of first row
		util.NewGameObject(&util.Point{X: sx - 2*s, Y: sy + 2*c}, tilt, assets.Tiles[436], util.NoCollider), // glass door
		util.NewGameObject(&util.Point{X: sx - 2*s + c, Y: sy + 2*c + s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 2*s + 2*c, Y: sy + 2*c + 2*s}, tilt, assets.Tiles[43], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 2*s + 3*c, Y: sy + 2*c + 3*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 2*s + 4*c, Y: sy + 2*c + 4*s}, tilt, assets.Tiles[45], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 2*s + 5*c, Y: sy + 2*c + 5*s}, tilt, assets.Tiles[45], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 2*s + 6*c, Y: sy + 2*c + 6*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 2*s + 7*c, Y: sy + 2*c + 7*s}, tilt, assets.Tiles[43], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 2*s + 8*c, Y: sy + 2*c + 8*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 2*s + 9*c, Y: sy + 2*c + 9*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 2*s + 10*c, Y: sy + 2*c + 10*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 2*s + 10*c, Y: sy + 2*c + 10*s}, tilt, assets.Tiles[437], util.NoCollider), // glass door

		util.NewGameObject(&util.Point{X: sx - 3*s, Y: sy + 3*c}, tilt, assets.Tiles[41], util.NoCollider),  // start of 2nd row
		util.NewGameObject(&util.Point{X: sx - 3*s, Y: sy + 3*c}, tilt, assets.Tiles[436], util.NoCollider), // glass door
		util.NewGameObject(&util.Point{X: sx - 3*s + c, Y: sy + 3*c + s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 3*s + 2*c, Y: sy + 3*c + 2*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 3*s + 3*c, Y: sy + 3*c + 3*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 3*s + 4*c, Y: sy + 3*c + 4*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 3*s + 5*c, Y: sy + 3*c + 5*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 3*s + 6*c, Y: sy + 3*c + 6*s}, tilt, assets.Tiles[43], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 3*s + 7*c, Y: sy + 3*c + 7*s}, tilt, assets.Tiles[45], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 3*s + 8*c, Y: sy + 3*c + 8*s}, tilt, assets.Tiles[45], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 3*s + 9*c, Y: sy + 3*c + 9*s}, tilt, assets.Tiles[45], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 3*s + 10*c, Y: sy + 3*c + 10*s}, tilt, assets.Tiles[45], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 3*s + 10*c, Y: sy + 3*c + 10*s}, tilt, assets.Tiles[437], util.NoCollider), // glass door

		util.NewGameObject(&util.Point{X: sx + 6*c - 1.5*s, Y: sy + 5*s + 1.5*c}, tilt, topHalfWoldenFloor, util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 7*c - 1.5*s, Y: sy + 6*s + 1.5*c}, tilt, topHalfWoldenFloor, util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 8*c - 1.5*s, Y: sy + 7*s + 1.5*c}, tilt, topHalfWoldenFloor, util.NoCollider),

		util.NewGameObject(&util.Point{X: sx - 4*s, Y: sy + 4*c}, tilt, assets.Tiles[43], util.NoCollider),  // start of 3rd row
		util.NewGameObject(&util.Point{X: sx - 4*s, Y: sy + 4*c}, tilt, assets.Tiles[436], util.NoCollider), // glass door
		util.NewGameObject(&util.Point{X: sx - 4*s + c, Y: sy + 4*c + s}, tilt, assets.Tiles[43], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 4*s + 2*c, Y: sy + 4*c + 2*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 4*s + 3*c, Y: sy + 4*c + 3*s}, tilt, assets.Tiles[45], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 4*s + 4*c, Y: sy + 4*c + 4*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 4*s + 5*c, Y: sy + 4*c + 5*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 4*s + 6*c, Y: sy + 4*c + 6*s}, tilt, assets.Tiles[45], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 4*s + 7*c, Y: sy + 4*c + 7*s}, tilt, assets.Tiles[43], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 4*s + 8*c, Y: sy + 4*c + 8*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 4*s + 9*c, Y: sy + 4*c + 9*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 4*s + 10*c, Y: sy + 4*c + 10*s}, tilt, assets.Tiles[41], util.NoCollider),

		util.NewGameObject(&util.Point{X: sx - 5*s + c, Y: sy + 5*c + s}, tilt, assets.Tiles[43], util.NoCollider), // start of 4th row
		util.NewGameObject(&util.Point{X: sx - 5*s + c, Y: sy + 5*c + s}, tilt, assets.Tiles[43], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 5*s + 2*c, Y: sy + 5*c + 2*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 5*s + 3*c, Y: sy + 5*c + 3*s}, tilt, assets.Tiles[43], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 5*s + 4*c, Y: sy + 5*c + 4*s}, tilt, assets.Tiles[43], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 5*s + 5*c, Y: sy + 5*c + 5*s}, tilt, assets.Tiles[45], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 5*s + 6*c, Y: sy + 5*c + 6*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 5*s + 7*c, Y: sy + 5*c + 7*s}, tilt, assets.Tiles[45], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 5*s + 8*c, Y: sy + 5*c + 8*s}, tilt, assets.Tiles[43], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 5*s + 9*c, Y: sy + 5*c + 9*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 5*s + 10*c, Y: sy + 5*c + 10*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 5*s + 10*c, Y: sy + 5*c + 10*s}, tilt, assets.Tiles[437], util.NoCollider), // glass door

		util.NewGameObject(&util.Point{X: sx - 6*s + c, Y: sy + 6*c + s}, tilt, assets.Tiles[41], util.NoCollider), // start of 5th row
		util.NewGameObject(&util.Point{X: sx - 6*s + c, Y: sy + 6*c + s}, tilt, assets.Tiles[43], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 6*s + 2*c, Y: sy + 6*c + 2*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 6*s + 3*c, Y: sy + 6*c + 3*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 6*s + 4*c, Y: sy + 6*c + 4*s}, tilt, assets.Tiles[45], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 6*s + 5*c, Y: sy + 6*c + 5*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 6*s + 6*c, Y: sy + 6*c + 6*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 6*s + 7*c, Y: sy + 6*c + 7*s}, tilt, assets.Tiles[41], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 6*s + 8*c, Y: sy + 6*c + 8*s}, tilt, assets.Tiles[43], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 6*s + 9*c, Y: sy + 6*c + 9*s}, tilt, assets.Tiles[45], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 6*s + 10*c, Y: sy + 6*c + 10*s}, tilt, assets.Tiles[43], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - 6*s + 10*c, Y: sy + 6*c + 10*s}, tilt, assets.Tiles[437], util.NoCollider), // glass door

		// kitchen area
		util.NewGameObject(&util.Point{X: sx - s + c, Y: sy + c + s}, tilt+math.Pi, assets.Tiles[323], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - s + 2*c, Y: sy + c + 2*s}, tilt+math.Pi, assets.Tiles[322], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - s + 3*c, Y: sy + c + 3*s}, tilt+math.Pi, assets.Tiles[320], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx - s + 3*c, Y: sy + c + 3*s}, tilt+math.Pi, assets.Tiles[267], util.NoCollider),

		// right middle corner
		util.NewGameObject(&util.Point{X: sx + 9*c - s, Y: sy + 9*s + c}, tilt, assets.Tiles[141], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 10*c - s, Y: sy + 10*s + c}, tilt, assets.Tiles[136], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 10*c, Y: sy + 10*s}, tilt, assets.Tiles[137], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 10*c + s, Y: sy + 10*s - c}, tilt, assets.Tiles[140], util.RectCollider),

		// upper right corner
		util.NewGameObject(&util.Point{X: sx + 4*c + 4*s, Y: sy + 4*s - 5*c}, tilt, assets.Tiles[141], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 5*c + 4*s, Y: sy + 5*s - 5*c}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 6*c + 4*s, Y: sy + 6*s - 5*c}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 7*c + 4*s, Y: sy + 7*s - 5*c}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 8*c + 4*s, Y: sy + 8*s - 5*c}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 9*c + 4*s, Y: sy + 9*s - 5*c}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 10*c + 4*s, Y: sy + 10*s - 5*c}, tilt, assets.Tiles[109], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 10*c + 3*s, Y: sy + 10*s - 4*c}, tilt, assets.Tiles[137], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx + 10*c + 2*s, Y: sy + 10*s - 3*c}, tilt, assets.Tiles[114], util.RectCollider),

		// wooden doors
		util.NewGameObject(&util.Point{X: sx + 6*c - s, Y: sy + 6*s + c}, tilt, assets.Tiles[467], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 7*c - s, Y: sy + 7*s + c}, tilt, assets.Tiles[467], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 8*c - s, Y: sy + 8*s + c}, tilt, assets.Tiles[467], util.NoCollider),

		util.NewGameObject(&util.Point{X: sx + 10*c + 2*s, Y: sy + 10*s - 2*c}, tilt, assets.Tiles[440], util.NoCollider),

		// pillars
		{
			Center:   &util.Point{X: sx + 5*c - 4*s, Y: sy + 5*s + 4*c},
			Rotation: tilt,
			Sprite:   assets.Tiles[170],
			Collider: util.NewRect(&util.Point{X: sx + 5*c - 4*s, Y: sy + 5*s + 4*c}, 48, 48, tilt),
		},
		util.NewGameObject(&util.Point{X: sx + 10*c - 5*s, Y: sy + 9*s + 4*c}, tilt, assets.Tiles[195], util.RectCollider),

		// bottom block
		util.NewGameObject(&util.Point{X: sx - 5*s, Y: sy + 5*c}, tilt, assets.Tiles[140], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 6*s, Y: sy + 6*c}, tilt, assets.Tiles[137], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 7*s, Y: sy + 7*c}, tilt, assets.Tiles[135], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 7*s + c, Y: sy + 7*c + s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 7*s + 2*c, Y: sy + 7*c + 2*s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 7*s + 3*c, Y: sy + 7*c + 3*s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 7*s + 4*c, Y: sy + 7*c + 4*s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 7*s + 5*c, Y: sy + 7*c + 5*s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 7*s + 6*c, Y: sy + 7*c + 6*s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 7*s + 7*c, Y: sy + 7*c + 7*s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 7*s + 8*c, Y: sy + 7*c + 8*s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 7*s + 9*c, Y: sy + 7*c + 9*s}, tilt, assets.Tiles[110], util.RectCollider),
		util.NewGameObject(&util.Point{X: sx - 7*s + 10*c, Y: sy + 7*c + 10*s}, tilt, assets.Tiles[113], util.RectCollider),

		// other objects without colliders

		util.NewGameObject(&util.Point{X: sx + 9*c + 4*s, Y: sy + 9*s - 4*c}, tilt+10.0*math.Pi/180.0, assets.Tiles[129], util.NoCollider),

		// other objects with colliders

		// small green couch on concrete floor
		{
			Center:   &util.Point{X: sx + 9*c, Y: sy + 9*s},
			Rotation: tilt - math.Pi/4,
			Sprite:   assets.Tiles[449],
			Collider: util.NewCircle(&util.Point{X: sx + 9*c, Y: sy + 9*s}, 26),
		},
		// green plant on concrete floor
		{
			Center:   &util.Point{X: sx + 5*c, Y: sy + 5*s},
			Rotation: tilt,
			Sprite:   assets.Tiles[133],
			Collider: util.NewCircle(&util.Point{X: sx + 5*c, Y: sy + 5*s}, 12),
		},
		// table on concrete floor
		{
			Center:   &util.Point{X: sx + 7*c + 2*s, Y: sy + 6*s - 2*c},
			Rotation: -tilt,
			Sprite:   assets.Tiles[510],
			Collider: util.NewRect(&util.Point{X: sx + 7*c + 2*s, Y: sy + 6*s - 2*c}, 24, 24, -tilt),
		},
		// chairs beside the wooden table below
		{ // top left
			Center:   &util.Point{X: sx + 2*c - 6*s, Y: sy + 2.5*c + 8*s},
			Rotation: -tilt - rotate180,
			Sprite:   assets.Tiles[528],
			Collider: util.NewRect(&util.Point{X: sx + 2*c - 5*s, Y: sy + 2.5*c + 6*s}, 8, 8, -tilt-rotate180),
		},
		{ // top right
			Center:   &util.Point{X: sx + 3*c - 7*s, Y: sy + 2.5*c + 7*s},
			Rotation: -tilt - rotate180,
			Sprite:   assets.Tiles[528],
			Collider: util.NewRect(&util.Point{X: sx + 3*c - 5*s, Y: sy + 2.5*c + 7*s}, 8, 8, -tilt-rotate180),
		},
		{ // bottom left
			Center:   &util.Point{X: sx + 2*c - 5*s, Y: sy + 3.5*c + 4*s},
			Rotation: -tilt,
			Sprite:   assets.Tiles[528],
			Collider: util.NewRect(&util.Point{X: sx + 2*c - 5*s, Y: sy + 3.5*c + 4*s}, 8, 8, -tilt),
		},
		{ // bottom right
			Center:   &util.Point{X: sx + 3*c - 5*s, Y: sy + 3.5*c + 8*s},
			Rotation: -8 * tilt,
			Sprite:   assets.Tiles[528],
			Collider: util.NewRect(&util.Point{X: sx + 3*c - 5*s, Y: sy + 3.5*c + 8*s}, 8, 8, -8*tilt),
		},
		// wooden table on wooden floor
		{
			Center:   &util.Point{X: sx + 2*c - 6*s, Y: sy + 3*c + 6*s},
			Rotation: -tilt,
			Sprite:   assets.Tiles[452],
			Collider: util.NewRect(&util.Point{X: sx + 2*c - 6*s, Y: sy + 3*c + 6*s}, 24, 24, -tilt),
		},
		{
			Center:   &util.Point{X: sx + 3*c - 6*s, Y: sy + 3*c + 5*s},
			Rotation: -tilt,
			Sprite:   assets.Tiles[454],
			Collider: util.NewRect(&util.Point{X: sx + 3*c - 6*s, Y: sy + 3*c + 5*s}, 24, 24, -tilt),
		},
		// stuff on the wooden table above
		util.NewGameObject(&util.Point{X: sx + 2.5*c - 6*s, Y: sy + 3*c + 5.5*s}, 0, assets.Tiles[239], util.NoCollider), // plant
		util.NewGameObject(&util.Point{X: sx + 3*c - 6*s, Y: sy + 3*c + 5*s}, 0, assets.Tiles[214], util.NoCollider),     // dishes
		// fridge on wooden floor
		util.NewGameObject(&util.Point{X: sx - 6*s + 2*c, Y: sy + 6*c + 2*s}, tilt, assets.Tiles[269], util.RectCollider),
		// lamp on wooden floor
		{
			Center:   &util.Point{X: sx + 4*c, Y: sy + 6*c},
			Rotation: 0,
			Sprite:   assets.Tiles[131],
			Collider: util.NewCircle(&util.Point{X: sx + 4*c, Y: sy + 6*c}, 4),
		},
		// carpet on wooden floor (coord is a bit fucked i know, could prolly be adjusted better)
		util.NewGameObject(&util.Point{X: sx + 7*c - 5*s, Y: sy + 5*s + 3*c}, -30*math.Pi/180.0, assets.Tiles[342], util.NoCollider),
		util.NewGameObject(&util.Point{X: sx + 8*c - 8*s, Y: sy + 14.8*s + 2*c}, -30*math.Pi/180.0, assets.Tiles[344], util.NoCollider),
		// glass table on wooden floor (duplicated to decrease opacity)
		{
			Center:   &util.Point{X: sx + 7*c, Y: sy + 4*s + 3*c},
			Rotation: 2 * tilt,
			Sprite:   assets.Tiles[479],
			Collider: util.NewRect(&util.Point{X: sx + 7*c, Y: sy + 4*s + 3*c}, 24, 24, 2*tilt),
		},
		util.NewGameObject(&util.Point{X: sx + 7*c, Y: sy + 4*s + 3*c}, 2*tilt, assets.Tiles[479], util.NoCollider),
		{
			Center:   &util.Point{X: sx + 8*c, Y: sy + 6*s + 3*c},
			Rotation: 2 * tilt,
			Sprite:   assets.Tiles[481],
			Collider: util.NewRect(&util.Point{X: sx + 8*c, Y: sy + 6*s + 3*c}, 24, 24, 2*tilt),
		},
		util.NewGameObject(&util.Point{X: sx + 8*c, Y: sy + 6*s + 3*c}, 2*tilt, assets.Tiles[481], util.NoCollider),
		// stuff on the glass table above
		util.NewGameObject(&util.Point{X: sx + 7*c + 2*s, Y: sy + 6*s + 3*c}, tilt, assets.Tiles[213], util.NoCollider),  // plate
		util.NewGameObject(&util.Point{X: sx + 7*c + 10*s, Y: sy + 6*s + 3*c}, tilt, assets.Tiles[240], util.NoCollider), // knife
		// green couch on wooden floor
		{
			Center:   &util.Point{X: sx + 7*c - 8*s, Y: sy + 2*s + 5*c},
			Rotation: tilt,
			Sprite:   assets.Tiles[446],
			Collider: util.NewRect(&util.Point{X: sx + 7*c - 8*s, Y: sy + 2*s + 5*c}, 24, 24, tilt),
		},
		{
			Center:   &util.Point{X: sx + 8*c - 8*s, Y: sy + 3*s + 5*c},
			Rotation: tilt,
			Sprite:   assets.Tiles[448],
			Collider: util.NewRect(&util.Point{X: sx + 8*c - 8*s, Y: sy + 3*s + 5*c}, 24, 24, tilt),
		},
		// green plant on wooden floor
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
