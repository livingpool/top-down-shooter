package bullet

import (
	"math"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/livingpool/top-down-shooter/game/assets"
	"github.com/livingpool/top-down-shooter/singleplayer/util"
)

type Bullet struct {
	ID     uuid.UUID
	Object *util.GameObject
}

func NewBullet(pos *util.Point, rotation float64) *Bullet {
	sprite := assets.Bullet

	return &Bullet{
		ID: uuid.New(),
		Object: &util.GameObject{
			Center:   pos,
			Rotation: rotation,
			Sprite:   sprite,
			Collider: util.NewCircle(pos, 4),
		},
	}
}

func (b *Bullet) Update() {
	speed := util.BulletSpeedPerSecond / float64(ebiten.TPS())

	b.Object.Center.X += math.Sin(b.Object.Rotation) * speed
	b.Object.Center.Y += math.Cos(b.Object.Rotation) * -speed
}

func (b *Bullet) Draw(screen *ebiten.Image, debugMode bool) {
	op := b.Object.CenterAndRotateImage()
	screen.DrawImage(b.Object.Sprite, op)

	if debugMode {
		b.Object.DrawDebugCircle(screen, 4, "")
	}
}
