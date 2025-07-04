package spawner

import (
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/livingpool/top-down-shooter/game/assets"
	"github.com/livingpool/top-down-shooter/singleplayer/util"
)

type ZombieSpawner struct {
	timer   *util.Timer
	zombies []*Zombie
}

// Every duration `d`, spawn `count` zombies
func NewZombieSpawner(d time.Duration, count int) *ZombieSpawner {
	return &ZombieSpawner{
		timer:   util.NewTimer(d),
		zombies: make([]*Zombie, count),
	}
}

func (zs *ZombieSpawner) Update() {
	zs.timer.Update()
	if zs.timer.IsReady() {
		zs.timer.Reset()
	}
}

func (zs *ZombieSpawner) Draw(screen *ebiten.Image) {
	for _, z := range zs.zombies {
		bounds := z.object.Sprite.Bounds()
		halfW := float64(bounds.Dx()) / 2
		halfH := float64(bounds.Dy()) / 2

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-halfW, -halfH)
		op.GeoM.Rotate(z.object.Rotation)
		op.GeoM.Translate(halfW, halfH)

		op.GeoM.Translate(z.object.X, z.object.Y)

		screen.DrawImage(z.object.Sprite, op)
	}
}

// Generates a random position at the edge of the screen (a ring).
// Targets
func randPosition(screenWidth, screenHeight float64, target util.Vector) util.Vector {
	// the distance from the center the zombie should spawn at — half the width
	r := screenWidth / 2.0

	// pick a random angle — 2π is 360° — so this returns 0° to 360°
	angle := rand.Float64() * 2 * math.Pi

	// figure out the spawn position by moving r pixels from the target at the chosen angle
	pos := util.Vector{
		X: target.X + math.Cos(angle)*r,
		Y: target.Y + math.Sin(angle)*r,
	}

	return pos
}

// Generates a random velocity
func randVelocity()

type Zombie struct {
	object util.GameObject
}

func NewZombie() *Zombie {
	sprite := assets.Zombie1HoldSprite

	return &Zombie{
		object: util.GameObject{
			Vector:   util.Vector{},
			Rotation: -util.FacingOffset,
			Sprite:   sprite,
		},
	}
}

func (z *Zombie) Update() {
}

func (z *Zombie) Collider() util.Rect {
	bounds := z.object.Sprite.Bounds()

	return util.NewRect(
		z.object.Vector.X,
		z.object.Vector.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}
