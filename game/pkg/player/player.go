package player

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/livingpool/top-down-shooter/game/assets"
	"github.com/livingpool/top-down-shooter/game/pkg/bullet"
	"github.com/livingpool/top-down-shooter/game/util"
)

const (
	translationPerSecond = 200
	shootCoolDown        = 500 * time.Millisecond
	gunPointOffset       = 20.0 * math.Pi / 180.0
	bulletSpawnOffset    = 30.0
)

type Player struct {
	Object        util.GameObject
	LastDelta     util.Vector // render rotation at the last frame to keep the facing position correctly
	HumanoidState util.HumanoidState
	Health        int
	Bullets       []*bullet.Bullet
	ShootCoolDown *util.Timer
	Ammo          int
	Inputs        []Input // local history of inputs
}

func NewPlayer(screenWidth, screenHeight float64) *Player {
	sprite := assets.ManBlueGunSprite

	pos := util.Vector{
		X: screenWidth / 2,
		Y: screenHeight / 2,
	}

	return &Player{
		Object: util.GameObject{
			Vector:   pos,
			Rotation: -util.FacingOffset,
			Sprite:   sprite,
		},
		LastDelta:     pos,
		HumanoidState: util.HumanoidStateStand,
		Bullets:       make([]*bullet.Bullet, 0),
		ShootCoolDown: util.NewTimer(shootCoolDown),
		Ammo:          0,
	}
}

func (p *Player) Update() {
	// move 200 pixels per second
	speed := float64(translationPerSecond / ebiten.TPS())

	var delta util.Vector

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		delta.X -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		delta.X += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		delta.Y -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		delta.Y += speed
	}

	// check for diagonal movement
	if delta.X != 0 && delta.Y != 0 {
		factor := speed / math.Sqrt(delta.X*delta.X+delta.Y*delta.Y)
		delta.X *= factor
		delta.Y *= factor
	}

	p.Object.Vector.X += delta.X
	p.Object.Vector.Y += delta.Y

	// update rotation
	if delta.X != 0 || delta.Y != 0 {
		p.Object.Rotation = math.Atan2(p.LastDelta.Y, p.LastDelta.X)
		p.LastDelta = delta
	}

	// update existing bullets
	for _, bullet := range p.Bullets {
		bullet.Update()
	}

	// shoot at specified intervals
	p.ShootCoolDown.Update()
	if p.ShootCoolDown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.ShootCoolDown.Reset()

		spawnPos := p.Object.CalcBulletSpawnPosition()

		bullet := bullet.NewBullet(spawnPos, p.Object.Rotation+util.FacingOffset)
		p.Bullets = append(p.Bullets, bullet)
	}
}

func (p *Player) Draw(screen *ebiten.Image, debugMode bool) {
	// draw player
	op := p.Object.CenterAndRotateImage()
	op.GeoM.Translate(p.Object.Vector.X, p.Object.Vector.Y)
	screen.DrawImage(p.Object.Sprite, op)

	// draw bullets
	for _, bullet := range p.Bullets {
		bullet.Draw(screen)
	}

	if debugMode {
		p.Object.Vector.DrawDebugCircle(screen, 32)
	}
}

func (p *Player) Collider() util.Rect {
	bounds := p.Object.Sprite.Bounds()

	return util.NewRect(
		p.Object.Vector.X,
		p.Object.Vector.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}

type Input struct {
	Seq       int
	TimeStamp int
	Inputs    []rune
}
