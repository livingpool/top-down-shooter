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
		bounds := z.Object.Sprite.Bounds()
		halfW := float64(bounds.Dx()) / 2
		halfH := float64(bounds.Dy()) / 2

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-halfW, -halfH)
		op.GeoM.Rotate(z.Object.Rotation)
		op.GeoM.Translate(halfW, halfH)

		op.GeoM.Translate(z.Object.X, z.Object.Y)

		screen.DrawImage(z.Object.Sprite, op)
	}
}

// Generates a random position at the edge of the screen (a ring).
// Targets
func randPosition(screenWidth, screenHeight float64, target util.Point) util.Point {
	// the distance from the center the zombie should spawn at — half the width
	r := screenWidth / 2.0

	// pick a random angle — 2π is 360° — so this returns 0° to 360°
	angle := rand.Float64() * 2 * math.Pi

	// figure out the spawn position by moving r pixels from the target at the chosen angle
	pos := util.Point{
		X: target.X + math.Cos(angle)*r,
		Y: target.Y + math.Sin(angle)*r,
	}

	return pos
}

// Generates a random velocity
func randVelocity()

type Zombie struct {
	Object util.GameObject
}

func NewZombie() *Zombie {
	sprite := assets.Zombie1HoldSprite

	return &Zombie{
		Object: util.GameObject{
			Point:    util.Point{},
			Rotation: -util.FacingOffset,
			Sprite:   sprite,
		},
	}
}

func (z *Zombie) Update() {
}

func (z *Zombie) Collider() util.Circle {
	bounds := z.Object.Sprite.Bounds()

	return util.NewCircle(
		z.Object.Point,
		max(float64(bounds.Dx())/2, float64(bounds.Dy())/2),
	)
}
