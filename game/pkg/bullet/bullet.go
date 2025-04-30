package bullet

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/livingpool/top-down-shooter/game/assets"
	"github.com/livingpool/top-down-shooter/game/util"
)

const (
	bulletSpeedPerSecond = 350.0
)

type Bullet struct {
	object util.GameObject
}

func NewBullet(pos util.Vector, rotation float64) *Bullet {
	sprite := assets.Bullet

	return &Bullet{
		object: util.GameObject{
			Vector:   pos,
			Rotation: rotation,
			Sprite:   sprite,
		},
	}
}

func (b *Bullet) Update() {
	speed := bulletSpeedPerSecond / float64(ebiten.TPS())

	b.object.Vector.X += math.Sin(b.object.Rotation) * speed
	b.object.Vector.Y += math.Cos(b.object.Rotation) * -speed
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	op := b.object.CenterAndRotateImage()

	op.GeoM.Scale(0.5, 0.5)
	op.GeoM.Translate(b.object.Vector.X, b.object.Vector.Y)

	screen.DrawImage(b.object.Sprite, op)
}

func (b *Bullet) Collider() util.Rect {
	bounds := b.object.Sprite.Bounds()

	return util.NewRect(
		b.object.Vector.X,
		b.object.Vector.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}
