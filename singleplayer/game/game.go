package game

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/livingpool/top-down-shooter/singleplayer/pkg/background"
	"github.com/livingpool/top-down-shooter/singleplayer/pkg/bullet"
	"github.com/livingpool/top-down-shooter/singleplayer/pkg/player"
	"github.com/livingpool/top-down-shooter/singleplayer/pkg/spawner"
	"github.com/livingpool/top-down-shooter/singleplayer/util"
)

// TODO: set global boundaries
type Game struct {
	DebugMode  bool
	Background *background.Background
	Player     *player.Player
	Camera     *util.Camera // camera follows the player's movements, but centered at {0, 0} initially
	Bullets    map[uuid.UUID]*bullet.Bullet
	Spawner    *spawner.ZombieSpawner
}

func NewGame(debugMode bool) *Game {
	logLevel := new(slog.LevelVar)
	if debugMode {
		logLevel.Set(slog.LevelInfo) // LevelDebug or LevelInfo
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: logLevel},
	)))

	camera := util.InitCamera()

	return &Game{
		DebugMode:  debugMode,
		Background: background.NewBackground(),
		Player:     player.NewPlayer("You", camera),
		Camera:     camera,
		Bullets:    make(map[uuid.UUID]*bullet.Bullet),
		Spawner:    spawner.NewZombieSpawner(5*time.Second, 3),
	}
}

func (g *Game) Update() error {
	newBullet := g.Player.Update(g.Camera)
	for _, b := range g.Bullets {
		b.Update()
	}
	if newBullet != nil {
		g.Bullets[uuid.New()] = newBullet
	}

	g.Spawner.Update(g.Player.Object.Center)

	g.ResolveCollisions()

	return nil
}

// Note that order determines the z-index
func (g *Game) Draw(screen *ebiten.Image) {
	g.Background.Draw(screen, 0, 0, g.DebugMode)

	g.Player.Draw(screen, g.DebugMode)
	for _, b := range g.Bullets {
		b.Draw(screen, g.DebugMode)
	}

	g.Spawner.Draw(screen)

	if g.DebugMode {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return util.ScreenWidth, util.ScreenHeight
}

func (g *Game) Reset() {
}
