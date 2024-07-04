// PLayer struct and methods
// for Asteroids written in Go using Ebitengine
// Author Paul Brace
// July 2024

package main

import (
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/paul63/vector2"
)

const (
	MaxThrust = 30
	GapTimer = 30
)

var (
	playerSprite = LoadImage("assets/player.png")
	playerSpriteThrust = LoadImage("assets/thrust.png")
	reloadTimer = 0
	reloadTime = 15
	rotationSpeed = math.Pi / float64(ebiten.TPS())
)

type Player struct {
	position vector2.Vector
	rotation float64
	sprite   *ebiten.Image
	width	 int
	height	 int
	loaded	bool			
	alive	bool
	thrust	float64
	direction	vector2.Vector
	hyperJumpTimer	int
	reverseTimer 	int
}

func NewPlayer() *Player {
	sprite := playerSprite
	pos := vector2.Vector{
		X: ScreenWidth/2,
		Y: ScreenHeight/2,
	}

	return &Player{
		position: pos,
		rotation: 0, 
		sprite:   sprite,
		width:	  sprite.Bounds().Dx(),
		height:   sprite.Bounds().Dy(),
		loaded:   true,
		alive:	  true,
		thrust:	  0,		
		hyperJumpTimer: 0,
		reverseTimer: 0,	
	}
}

func (p *Player) ScreenPos() ScreenPos {
	return ScreenPos{
		X: 		p.position.X,
		Y: 		p.position.Y,
		Width: 	p.width,
		Height: p.height,
	}

}

func (p *Player) LaunchMissile() {
	NewMissile(p.position, p.rotation)
	p.loaded = false
	reloadTimer = reloadTime
}

func (p *Player) Update() {
	// check if need to move player
	if p.thrust > 0 {
		p.position.X += p.direction.X * p.thrust / 5
		p.position.Y += p.direction.Y * p.thrust / 5
		p.thrust -= 1
		if p.thrust == 0 {
			p.sprite = playerSprite
		}
		if p.position.X > ScreenWidth {
			p.position.X = 0
		} else {
			if p.position.X < 0 {
				p.position.X = ScreenWidth
			}
		}
		if p.position.Y > ScreenHeight{
			p.position.Y = 0
		} else {
			if p.position.Y < 0 {
				p.position.Y = ScreenHeight
			}
		}
	}

	// Update gap timers
	if p.hyperJumpTimer > 0 {
		p.hyperJumpTimer -= 1
	}
	if p.reverseTimer > 0 {
		p.reverseTimer -= 1
	}
	if reloadTimer > 0{
		reloadTimer -= 1
		if reloadTimer == 0 {
			p.loaded = true
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		// rotate left
		p.rotation -= rotationSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		// rotate right
		p.rotation += rotationSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && p.reverseTimer <= 0 {
		// reverse direction
		p.rotation += 1.5708 * 2
		p.reverseTimer = GapTimer
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) && p.loaded {
		// fire missile
		p.LaunchMissile()
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) && p.loaded {
        // Set player to point in direction of mouse cursor and fire missile
		x, y := ebiten.CursorPosition()
		p.rotation = p.position.PointTowards(vector2.Vector{X: float64(x), Y: float64(y)})
		p.LaunchMissile()
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton2) || ebiten.IsKeyPressed(ebiten.KeyArrowUp){
		// Set player movement in progress
		p.thrust = MaxThrust
		p.sprite = playerSpriteThrust
		// Calculate a target so flies in direction ship pointing
		p.direction = vector2.Vector{
			X: math.Sin(p.rotation),
			Y: math.Cos(p.rotation) * -1,
		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton1) || ebiten.IsKeyPressed(ebiten.KeyH){
		// Perform hyperjump
		if p.hyperJumpTimer <= 0 {
			p.hyperJumpTimer = GapTimer
			p.position.X = float64(rand.IntN(ScreenWidth - 80) + 40)
			p.position.Y = float64(rand.IntN(ScreenHeight - 80) + 40)
		}
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	if p.alive{
		DrawImage(screen, p.sprite, p.position.X, p.position.Y, p.rotation)
	}
}

func (p *Player) Hit(){
	p.alive = false
	NewExplosion(p.position.X, p.position.Y, 75, 
		color.RGBA{255, 0, 0, 100}, 0.025)
}

func (p *Player) Reset() {
	p.position.X = ScreenWidth/2
	p.position.Y = ScreenHeight/2
	p.rotation = 0
	p.loaded = true
	p.alive = true
	p.thrust = 0
}
