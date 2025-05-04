package util

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GameObject struct {
	Vector           // center coord of object
	Rotation float64 // where the object is facing
	Sprite   *ebiten.Image
}

// Return a *ebiten.DrawImageOptions where the sprite is centered at (0,0) and rotated.
func (obj *GameObject) CenterAndRotateImage() *ebiten.DrawImageOptions {
	bounds := obj.Sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(obj.Rotation)

	return op
}

// Calculates the bullet's spawn position relative to the object.
// Note that we need to offset bullet's rotation due to object's facing direction.
// Maybe I would extract that out in the future.
func (obj *GameObject) CalcBulletSpawnPosition() Vector {
	spawnRotation := obj.Rotation + GunPointOffset
	spawnPos := Vector{
		X: obj.Vector.X + math.Cos(spawnRotation)*BulletSpawnOffset,
		Y: obj.Vector.Y + math.Sin(spawnRotation)*BulletSpawnOffset,
	}
	return spawnPos
}

func (obj *GameObject) DrawDebugCircle(radius float32, debugText string) {
	obj.Vector.DrawDebugCircle(obj.Sprite, radius)
	if debugText != "" {
		ebitenutil.DebugPrint(obj.Sprite, debugText)
	}
}

type Vector struct {
	X float64
	Y float64
}

func (v Vector) Normalize() Vector {
	magnitude := math.Sqrt(v.X*v.X + v.Y*v.Y)
	return Vector{v.X / magnitude, v.Y / magnitude}
}

func (v Vector) DrawDebugCircle(screen *ebiten.Image, radius float32) {
	c := color.RGBA{R: 74, G: 246, B: 38, A: 1}
	vector.StrokeCircle(screen, float32(v.X), float32(v.Y), radius, 1, c, true)
}

type Rect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

func NewRect(x, y, width, height float64) Rect {
	return Rect{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

func (r Rect) MaxX() float64 {
	return r.X + r.Width
}

func (r Rect) MaxY() float64 {
	return r.Y + r.Height
}

func (r Rect) Intersects(other Rect) bool {
	return r.X <= other.MaxX() &&
		other.X <= r.MaxX() &&
		r.Y <= other.MaxY() &&
		other.Y <= r.MaxY()
}
