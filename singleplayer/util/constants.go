package util

import (
	"math"
	"time"
)

// Game settings
const (
	ScreenWidth         = 800
	ScreenHeight        = 600
	ServerPhysicsPeriod = 15 * time.Millisecond
	ServerUpdatePeriod  = 45 * time.Millisecond
)

// Position offsets
const (
	FacingOffset      = 90.0 * math.Pi / 180.0
	GunPointOffset    = 20.0 * math.Pi / 180.0
	BulletSpawnOffset = 30.0
)

// Initial player states
const (
	InitialPlayerHealth   = 5
	InitialPlayerAmmo     = 10
	InitialPlayerX        = ScreenWidth / 2
	InitialPlayerY        = ScreenHeight/2 + 50
	InitialPlayerRotation = -FacingOffset
)

// All the different states a humanoid can be in
type HumanoidState int

const (
	HumanoidStateGun = iota
	HumanoidStateHold
	HumanoidStateMachine
	HumanoidStateReload
	HumanoidStateSilencer
	HumanoidStateStand
)

// Player settings
const (
	PlayerSpeedPerSecond = 200 // move x pixels per second
	PlayerShootCoolDown  = 500 * time.Millisecond
)

// Bullet settings
const (
	BulletSpeedPerSecond = 350.0
)

// Zombie spawner settings
const (
	ZombieMaxSpeedPerSecond = 300
	ZombieMinSpeedPerSecond = 100
)
