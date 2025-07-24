package util

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GameObject struct {
	Center   *Point        // center coord of the object
	Rotation float64       // where the object is facing
	Sprite   *ebiten.Image // sprite of the object
	Collider               // used to check collisions; can be nil
}

// NewGameObject creates a new GameObject with default collision zone,
// for something like bullet, make the struct yourself.
func NewGameObject(center *Point, rotation float64, sprite *ebiten.Image, colliderType ColliderType) *GameObject {
	obj := &GameObject{
		Center:   center,
		Rotation: rotation,
		Sprite:   sprite,
	}

	switch colliderType { // i hard coded the dimensions bcos all my sprites are of the same dimensions
	case RectCollider:
		obj.Collider = NewRect(center, 64, 64, rotation)
	case CircleCollider:
		obj.Collider = NewCircle(center, 32)
	default:
		log.Fatal("unsupported collider type")
	}

	return obj
}

// obj.CenterAndRotateImage returns a *ebiten.DrawImageOptions where the sprite is adjusted according to obj's center and rotation
func (obj GameObject) CenterAndRotateImage() *ebiten.DrawImageOptions {
	bounds := obj.Sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(obj.Rotation)

	op.GeoM.Translate(obj.Center.X, obj.Center.Y)

	return op
}

// Calculates the bullet's spawn position relative to the object.
// Note that we need to offset bullet's rotation due to object's facing direction.
// Maybe I would extract that out in the future.
func (obj GameObject) CalcBulletSpawnPosition() Point {
	spawnRotation := obj.Rotation + GunPointOffset
	spawnPos := Point{
		X: obj.Center.X + math.Cos(spawnRotation)*BulletSpawnOffset,
		Y: obj.Center.Y + math.Sin(spawnRotation)*BulletSpawnOffset,
	}
	return spawnPos
}

func (obj GameObject) DrawDebugCircle(screen *ebiten.Image, radius float32, debugText string) {
	obj.Center.DrawDebugCircle(screen, radius)
	if debugText != "" {
		ebitenutil.DebugPrint(obj.Sprite, debugText)
	}
}

func (obj GameObject) DrawDebugRect(screen *ebiten.Image, dimX, dimY float32, debugText string) {
	obj.Center.DrawDebugRect(screen, dimX, dimY)
	if debugText != "" {
		ebitenutil.DebugPrint(obj.Sprite, debugText)
	}
}

func (obj GameObject) Collide(other GameObject) (Vector, bool) {
	return obj.Collider.Collide(other.Collider)
}

func (obj GameObject) GetColliderType() ColliderType {
	switch obj.Collider.(type) {
	case Rect:
		return RectCollider
	case Circle:
		return CircleCollider
	default:
		log.Fatal("unrecognized collider type")
		return -1
	}
}

type Point struct {
	X float64
	Y float64
}

func (p *Point) Add(v Vector) {
	p.X += v.X
	p.Y += v.Y
}

func (p *Point) Sub(v Vector) {
	p.X -= v.X
	p.Y -= v.Y
}

func (p Point) Distance(other Point) float64 {
	xSquared := math.Pow(math.Abs(p.X-other.X), 2)
	ySquared := math.Pow(math.Abs(p.Y-other.Y), 2)
	return math.Sqrt(xSquared + ySquared)
}

func (p Point) ManhattanDistance(other Point) float64 {
	return math.Abs(p.X-other.X) + math.Abs(p.Y-other.Y)
}

func (p Point) DrawDebugCircle(screen *ebiten.Image, radius float32) {
	c := color.RGBA{R: 74, G: 246, B: 38, A: 1}
	vector.StrokeCircle(screen, float32(p.X), float32(p.Y), radius, 1, c, true)
}

// cant draw a rectangle with rotation i think
func (p Point) DrawDebugRect(screen *ebiten.Image, dimX, dimY float32) {
	c := color.RGBA{R: 74, G: 246, B: 38, A: 1}
	originX := float32(p.X) - dimX/2
	originY := float32(p.Y) - dimY/2
	vector.StrokeRect(screen, originX, originY, dimX, dimY, 1, c, true)
}

// p.Vector returns a Vector p -> other
func (p Point) Vector(other Point) Vector {
	return Vector{other.X - p.X, other.Y - p.Y}
}

type Vector struct {
	X float64
	Y float64
}

func (v *Vector) Add(other Vector) {
	v.X += other.X
	v.Y += other.Y
}

func (v Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vector) ReverseDirection() Vector {
	return Vector{-v.X, -v.Y}
}

func (v Vector) Normalize() Vector {
	magnitude := v.Length()
	if magnitude == 0 {
		log.Fatalf("vector %v has length = 0, returning it as is...", v)
	}
	return Vector{v.X / magnitude, v.Y / magnitude}
}

func (v Vector) Scale(n float64) Vector {
	if n == 0 {
		log.Fatal("scaling to zero is not allowed")
	}
	return Vector{v.X * n, v.Y * n}
}

func (v Vector) GetPerpendicularVector() Vector {
	return Vector{-v.Y, v.X}
}

func (v Vector) InnerProduct(other Vector) float64 {
	return v.X*other.X + v.Y*other.Y
}
