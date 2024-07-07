// Missile struct and methods
// for Asteroids written in Go using Ebitengine
// Author Paul Brace
// July 2024

package main

import (
	"math"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/paul63/vector2"
)

var (
	missileSprite = LoadImage("assets/bullet.png")
	missiles [] *Missile
)
func UpdateAllMissiles(){
	for _, m := range missiles {
		m.Update()
	} 
}

func DrawAllMissiles(screen *ebiten.Image){
	for _, m := range missiles {
		m.Draw(screen)
	} 
}

func ClearDoneMissiles(){
	for i, m := range(missiles){
		if m.done{
			missiles = append(missiles[:i], missiles[i+1:]...)
			break
		}		
	}
}

func ClearAllMissiles(){
	missiles = nil
}

type Missile struct {
	GameSprite
}

// Create missile at x, y, rotated by angle and traveling toward target
func NewMissile(pos vector2.Vector, angle float64) *Missile {
	speed := 5.0
	sprite := missileSprite

	// Calculate a target so flies in direction ship pointing
	target := vector2.Vector{
		X: pos.X + math.Sin(angle),
		Y: pos.Y + math.Cos(angle) * -1,
	}
	direction := vector2.Vector{
					X: target.X - pos.X,
					Y: target.Y - pos.Y,
				}
	direction.Normalize()

	movement := vector2.Vector{
					X: direction.X * speed,
					Y: direction.Y * speed,
	}

	// move forward so appears at front of ship
	pos.X += movement.X * 4
	pos.Y += movement.Y * 4

	gameSprite := NewGameSprite(sprite, pos, movement, angle)

	//create missile
	missile := Missile{
		GameSprite: gameSprite,
	}

	// Add to list of missiles
	missiles = append(missiles, &missile)

	return &missile
}

func (m *Missile) ScreenPos() ScreenPos {
	return ScreenPos{
		X: 		m.position.X,
		Y: 		m.position.Y,
		Width: 	m.width,
		Height: m.height,
	}

}

func (m *Missile) Update() {
	if !m.done {
		m.position.Add(m.movement)
		m.done =  m.position.X < -50 || m.position.X > ScreenWidth + 50 ||
				m.position.Y < -50 || m.position.Y > ScreenHeight + 50
	}
}

func (m *Missile) Draw(screen *ebiten.Image) {
	m.DrawImage(screen)
}
