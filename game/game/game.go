package game

import (
	"time"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/livingpool/top-down-shooter/game/pkg/bullet"
	"github.com/livingpool/top-down-shooter/game/pkg/player"
	"github.com/livingpool/top-down-shooter/game/util"
)

// The main game class. This gets created on both server and client.
// Server creates one game instance for each game that is hosted,
// and client creates one for itself to play the game.
// TODO: set boundaries
// TODO: upon game creation, client need to construct the Players map
type Game struct {
	ID        uuid.UUID
	DebugMode bool
	IsServer  bool // store a flag to determine if this instance is a server or client

	// not yet used
	PhysicsDelta          int
	PhysicsLastUpdateTime time.Time // in ms
	LocalTimeElapsed      int       // in seconds
	LocalDelta            int
	LocalLastFrameTime    time.Time // in ms

	Players map[uuid.UUID]*player.Player
	Bullets map[uuid.UUID]*bullet.Bullet
}

func NewGame(isServer bool) *Game {
	return &Game{
		ID:        uuid.New(),
		DebugMode: true,
		IsServer:  isServer,
		Players:   make(map[uuid.UUID]*player.Player),
		Bullets:   make(map[uuid.UUID]*bullet.Bullet),
	}
}

// The physics update loop
func (g *Game) Update() error {
	for _, p := range g.Players {
		p.Update()
	}
	for _, b := range g.Bullets {
		b.Update()
	}
	g.checkCollisions()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.Players {
		p.Draw(screen, g.DebugMode)
	}
	for _, b := range g.Bullets {
		b.Draw(screen, g.DebugMode)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return util.ScreenWidth, util.ScreenHeight
}

// TODO: randomize players' spawn positions; at both initial spawn and reset
func (g *Game) Reset() {
	for i := range g.Players {
		g.Players[i] = player.NewPlayer("todo")
	}
}
