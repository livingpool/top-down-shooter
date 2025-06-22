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

func (obj *GameObject) DrawDebugCircle(screen *ebiten.Image, radius float32, debugText string) {
	obj.Vector.DrawDebugCircle(screen, radius)
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

func (v Vector) InnerProduct(other Vector) float64 {
	return v.X*other.X + v.Y*other.Y
}

func (v Vector) DrawDebugCircle(screen *ebiten.Image, radius float32) {
	c := color.RGBA{R: 74, G: 246, B: 38, A: 1}
	vector.StrokeCircle(screen, float32(v.X), float32(v.Y), radius, 1, c, true)
}

// There are currently two Collider objects: Rect and Circle
// https://stackoverflow.com/questions/401847/circle-rectangle-collision-detection-intersection

type Rect struct {
	Center Vector
	MaxX   Vector
	MaxY   Vector
}

func NewRect(center, maxX, maxY Vector) Rect {
	v := Vector{center.X - maxX.X, center.Y - maxX.Y}
	v2 := Vector{center.X - maxY.X, center.Y - maxY.Y}
	if v.InnerProduct(v2) != 0 {
		log.Fatalf("failed to create rect: given points do not form a rect")
	}

	return Rect{
		Center: center,
		MaxX:   maxX,
		MaxY:   maxY,
	}
}

type Circle struct {
	Vector
	Radius float64
}

func NewCircle(x, y, radius float64) Circle {
	return Circle{
		Vector: Vector{x, y},
		Radius: radius,
	}
}

func IntersectRectAndCircle(r Rect, c Circle) bool {
}

func IntersectRectAndRect(r1, r2 Rect) bool {
}

func IntersectCircleAndCircle(c1, c2 Circle) bool {
}
