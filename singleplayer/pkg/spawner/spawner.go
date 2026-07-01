package spawner

import (
	"log/slog"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/livingpool/top-down-shooter/game/assets"
	"github.com/livingpool/top-down-shooter/singleplayer/util"
)

type ZombieSpawner struct {
	spawnDuration time.Duration
	spawnCount    int
	timer         *util.Timer
	zombies       []*Zombie
}

// Every duration `d`, spawn `count` zombies
func NewZombieSpawner(d time.Duration, count int) *ZombieSpawner {
	return &ZombieSpawner{
		spawnDuration: d,
		spawnCount:    count,
		timer:         util.NewTimer(d),
		zombies:       make([]*Zombie, 0, count),
	}
}

func (zs *ZombieSpawner) Update(target *util.Point) {
	zs.timer.Update()
	if zs.timer.IsReady() {
		for range zs.spawnCount {
			pos := randPosition(util.ScreenWidth, util.ScreenHeight, util.Point{X: 0, Y: 0})
			zombie := NewZombie(&pos, 0, float64(randSpeed()), 1, assets.Zombie1StandSprite, target)
			zs.zombies = append(zs.zombies, zombie)
			slog.Info("Zombie spawned", "pos", pos)
		}

		zs.timer.Reset()
	}

	for _, z := range zs.zombies {
		z.Update(target)
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

		op.GeoM.Translate(z.Object.Center.X, z.Object.Center.Y)

		screen.DrawImage(z.Object.Sprite, op)
	}
}

// Generates a random position at the edge of the screen (a ring).
func randPosition(screenWidth, screenHeight float64, target util.Point) util.Point {
	// the distance from the center the zombie should spawn at — half the width
	r := max(screenWidth, screenHeight) / 2.0

	// pick a random angle — 2π is 360° — so this returns 0° to 360°
	angle := rand.Float64() * 2 * math.Pi

	// figure out the spawn position by moving r pixels from the target at the chosen angle
	pos := util.Point{
		X: target.X + math.Cos(angle)*r,
		Y: target.Y + math.Sin(angle)*r,
	}

	return pos
}

type Zombie struct {
	Object   *util.GameObject
	Health   int
	Velocity float64
	Target   *util.Point // for now this tracks the player's position
}

func NewZombie(pos *util.Point, rot, velocity float64, health int, sprite *ebiten.Image, target *util.Point) *Zombie {
	return &Zombie{
		Object:   util.NewGameObject(pos, rot, sprite, util.CircleCollider),
		Health:   health,
		Velocity: velocity,
		Target:   target,
	}
}

// calc zombie's rotation wrt to the player's position
func (z *Zombie) Update(target *util.Point) {
	dx := target.X - z.Object.Center.X
	dy := target.Y - z.Object.Center.Y
	z.Object.Rotation = math.Atan2(dy, dx)
}

// Generates a random speed
func randSpeed() int {
	max, min := util.ZombieMaxSpeedPerSecond, util.ZombieMinSpeedPerSecond
	speed := rand.Intn(max-min) + min

	return speed
}
