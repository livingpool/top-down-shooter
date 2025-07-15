package player

import (
	"log/slog"
	"math"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/livingpool/top-down-shooter/game/assets"
	"github.com/livingpool/top-down-shooter/singleplayer/pkg/bullet"
	"github.com/livingpool/top-down-shooter/singleplayer/util"
)

type Player struct {
	ID            uuid.UUID
	Name          string
	Conn          *websocket.Conn
	Object        *util.GameObject
	LastDelta     util.Point // render rotation at the last frame to keep the facing position correctly
	HumanoidState util.HumanoidState
	Health        int
	ShootCoolDown *util.Timer
	Ammo          int
}

func NewPlayer(name string) *Player {
	sprite := assets.ManBlueGunSprite

	pos := util.Point{
		X: util.InitialPlayerX,
		Y: util.InitialPlayerY,
	}

	return &Player{
		ID:            uuid.New(),
		Name:          name,
		Object:        util.NewGameObject(&pos, util.InitialPlayerRotation, sprite, util.CircleCollider),
		LastDelta:     pos,
		HumanoidState: util.HumanoidStateStand,
		Health:        util.InitialPlayerHealth,
		ShootCoolDown: util.NewTimer(util.PlayerShootCoolDown),
		Ammo:          util.InitialPlayerAmmo,
	}
}

// Player.Update() updates the player and returns a new Bullet (can be nil).
func (p *Player) Update() *bullet.Bullet {
	// move 200 pixels per second
	speed := float64(util.PlayerSpeedPerSecond / ebiten.TPS())

	var b *bullet.Bullet

	p.ShootCoolDown.Update()

	var delta util.Point

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

	p.Object.Center.X += delta.X
	p.Object.Center.Y += delta.Y

	// update rotation
	if delta.X != 0 || delta.Y != 0 {
		newRotation := math.Atan2(p.LastDelta.Y, p.LastDelta.X)
		slog.Debug("rotation updated", "old", p.Object.Rotation, "new", newRotation)

		p.Object.Rotation = newRotation
		p.LastDelta = delta
	}

	// constrain shooting at fixed intervals
	if p.ShootCoolDown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.ShootCoolDown.Reset()

		spawnPos := p.Object.CalcBulletSpawnPosition()
		b = bullet.NewBullet(&spawnPos, p.Object.Rotation+util.FacingOffset)

		slog.Debug("new bullet", "pos", spawnPos)
	}

	return b
}

func (p *Player) Draw(screen *ebiten.Image, debugMode bool) {
	op := p.Object.CenterAndRotateImage()

	screen.DrawImage(p.Object.Sprite, op)

	if debugMode {
		p.Object.DrawDebugCircle(screen, 32, p.Name)
	}
}
