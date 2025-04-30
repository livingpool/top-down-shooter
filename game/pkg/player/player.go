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
	object        util.GameObject
	lastDelta     util.Vector // render rotation at the last frame to keep the facing position correctly
	state         util.HumanoidState
	health        int
	Bullets       []*bullet.Bullet
	shootCoolDown *util.Timer
	ammo          int
}

func NewPlayer(screenWidth, screenHeight float64) *Player {
	sprite := assets.ManBlueGunSprite

	pos := util.Vector{
		X: screenWidth / 2,
		Y: screenHeight / 2,
	}

	return &Player{
		object: util.GameObject{
			Vector:   pos,
			Rotation: -util.FacingOffset,
			Sprite:   sprite,
		},
		lastDelta:     pos,
		state:         util.HumanoidStateStand,
		Bullets:       make([]*bullet.Bullet, 0),
		shootCoolDown: util.NewTimer(shootCoolDown),
		ammo:          0,
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

	p.object.Vector.X += delta.X
	p.object.Vector.Y += delta.Y

	// update rotation
	if delta.X != 0 || delta.Y != 0 {
		p.object.Rotation = math.Atan2(p.lastDelta.Y, p.lastDelta.X)
		p.lastDelta = delta
	}

	// update existing bullets
	for _, bullet := range p.Bullets {
		bullet.Update()
	}

	// shoot at specified intervals
	p.shootCoolDown.Update()
	if p.shootCoolDown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.shootCoolDown.Reset()

		spawnPos := p.object.CalcBulletSpawnPosition()

		bullet := bullet.NewBullet(spawnPos, p.object.Rotation+util.FacingOffset)
		p.Bullets = append(p.Bullets, bullet)
	}
}

func (p *Player) Draw(screen *ebiten.Image, debugMode bool) {
	// draw player
	op := p.object.CenterAndRotateImage()
	op.GeoM.Translate(p.object.Vector.X, p.object.Vector.Y)
	screen.DrawImage(p.object.Sprite, op)

	// draw bullets
	for _, bullet := range p.Bullets {
		bullet.Draw(screen)
	}

	if debugMode {
		p.object.Vector.DrawDebugCircle(screen, 32)
	}
}

func (p *Player) Collider() util.Rect {
	bounds := p.object.Sprite.Bounds()

	return util.NewRect(
		p.object.Vector.X,
		p.object.Vector.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}
