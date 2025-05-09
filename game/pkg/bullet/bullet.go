package bullet

import (
	"math"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/livingpool/top-down-shooter/game/assets"
	"github.com/livingpool/top-down-shooter/game/util"
)

type Bullet struct {
	ID     uuid.UUID
	Object util.GameObject
}

func NewBullet(pos util.Vector, rotation float64) *Bullet {
	sprite := assets.Bullet

	return &Bullet{
		ID: uuid.New(),
		Object: util.GameObject{
			Vector:   pos,
			Rotation: rotation,
			Sprite:   sprite,
		},
	}
}

func (b *Bullet) Update() {
	speed := util.BulletSpeedPerSecond / float64(ebiten.TPS())

	b.Object.Vector.X += math.Sin(b.Object.Rotation) * speed
	b.Object.Vector.Y += math.Cos(b.Object.Rotation) * -speed
}

func (b *Bullet) Draw(screen *ebiten.Image, debugMode bool) {
	op := b.Object.CenterAndRotateImage()

	op.GeoM.Scale(0.5, 0.5)
	op.GeoM.Translate(b.Object.Vector.X, b.Object.Vector.Y)

	screen.DrawImage(b.Object.Sprite, op)

	if debugMode {
		b.Object.DrawDebugCircle(screen, 4, "")
	}
}

func (b *Bullet) Collider() util.Rect {
	bounds := b.Object.Sprite.Bounds()

	return util.NewRect(
		b.Object.Vector.X,
		b.Object.Vector.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}
