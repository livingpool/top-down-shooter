package player

import (
	"math"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/livingpool/top-down-shooter/game/assets"
	"github.com/livingpool/top-down-shooter/game/pkg/bullet"
	"github.com/livingpool/top-down-shooter/game/util"
)

type Player struct {
	ID            uuid.UUID
	Name          string
	Conn          *websocket.Conn
	Object        util.GameObject
	LastDelta     util.Vector // render rotation at the last frame to keep the facing position correctly
	HumanoidState util.HumanoidState
	Health        int
	ShootCoolDown *util.Timer
	Ammo          int
	ClientUpdates []util.ClientUpdate // local history of inputs
	LastInputSeq  int
}

func NewPlayer(name string) *Player {
	sprite := assets.ManBlueGunSprite

	pos := util.Vector{
		X: util.InitialPlayerX,
		Y: util.InitialPlayerY,
	}

	return &Player{
		ID:   uuid.New(),
		Name: name,
		Object: util.GameObject{
			Vector:   pos,
			Rotation: util.InitialPlayerRotation,
			Sprite:   sprite,
		},
		LastDelta:     pos,
		HumanoidState: util.HumanoidStateStand,
		Health:        util.InitialPlayerHealth,
		ShootCoolDown: util.NewTimer(util.PlayerShootCoolDown),
		Ammo:          util.InitialPlayerAmmo,
		ClientUpdates: make([]util.ClientUpdate, 0),
		LastInputSeq:  0,
	}
}

// Player.Update() updates the player and returns a new Bullet (can be nil).
// Note that ShootCoolDown must > physics update period, or multiple bullets may be created.
// TODO: race condition at the slice ClientUpdates?
func (p *Player) Update() *bullet.Bullet {
	// move 200 pixels per second
	speed := float64(util.PlayerSpeedPerSecond / ebiten.TPS())

	var b *bullet.Bullet

	p.ShootCoolDown.Update()

	for _, msg := range p.ClientUpdates {
		// skip inputs we have already simulated locally
		if msg.Seq <= p.LastInputSeq {
			continue
		}

		var delta util.Vector

		input := msg.Keys
		if input.A {
			delta.X -= speed
		}
		if input.D {
			delta.X += speed
		}
		if input.W {
			delta.Y -= speed
		}
		if input.S {
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

		// constrain shooting at fixed intervals
		if p.ShootCoolDown.IsReady() && input.Space {
			p.ShootCoolDown.Reset()

			spawnPos := p.Object.CalcBulletSpawnPosition()
			b = bullet.NewBullet(spawnPos, p.Object.Rotation+util.FacingOffset)
		}

		p.LastInputSeq = msg.Seq
	}

	// TODO: clear array

	return b
}

func (p *Player) Draw(screen *ebiten.Image, debugMode bool) {
	op := p.Object.CenterAndRotateImage()
	op.GeoM.Translate(p.Object.Vector.X, p.Object.Vector.Y)
	screen.DrawImage(p.Object.Sprite, op)

	if debugMode {
		p.Object.DrawDebugCircle(screen, 32, "")
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
