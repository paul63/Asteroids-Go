// Meteor struct and methods
// for Asteroids written in Go using Ebitengine
// Author Paul Brace
// July 2024

package main

import (
	"image/color"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/paul63/vector2"
)

var (
	meteors [] *Meteor
	meteorSprites = [] *ebiten.Image {LoadImage("assets/asteroidTiny.png"),
 									LoadImage("assets/asteroidSmall.png"),
 									LoadImage("assets/asteroidMed.png"),
									LoadImage("assets/asteroidLarge.png")}
	scores = [] int {100, 75, 50, 25}
	expColor color.Color = color.RGBA{255, 125, 125, 100}
)

func UpdateAllMeteors(){
	for _, m := range meteors {
		m.Update()
	} 
}

func DrawAllMeteors(screen *ebiten.Image){
	for _, m := range meteors {
		m.Draw(screen)
	} 
}

func ClearDoneMeteors(){
	for i, m := range(meteors){
		if m.done{
			// Clear first done found - will eventually clear all done as called every frame
			meteors = append(meteors[:i], meteors[i+1:]...)
			break
		}		
	}
}

func ClearAllMeteors(){
	meteors = nil
}

type Meteor struct {
	position 		vector2.Vector
	movement		vector2.Vector
	rotationSpeed	float64
	angle			float64
	size			int
	sprite   		*ebiten.Image
	width 			int
	height			int
	done			bool
}

// To create a new Meteor and add to list
// Player position used as destination of meteor
func NewMeteor(player *Player) *Meteor {
	size := rand.IntN(4)
	sprite := meteorSprites[size]

	// set destination to player position
	target := vector2.Vector{
		X: player.position.X,
		Y: player.position.Y,
	}
	// select a screen edge to enter from
	edge := rand.IntN(4)
	var (
		x int
		y int
	)
	switch edge{
	case 0:
		x = -10
		y = rand.IntN(ScreenWidth)
	case 1:
		x = ScreenWidth + 10
		y = rand.IntN(ScreenHeight)
	case 2:
		y = - 10
		x = rand.IntN(ScreenWidth)
	default:
		y = ScreenHeight + 10
		x = rand.IntN(ScreenHeight)
	}
	pos := vector2.Vector{
		X: float64(x),
		Y: float64(y),
	}

	// Randomized velocity
	velocity := 0.25 + rand.Float64()*1.5

	// Direction is the target minus the current position
	direction := vector2.Vector{
		X: (target.X - pos.X),
		Y: (target.Y - pos.Y),
	}

	// remove distance from Vector
	direction.Normalize()

	movement := vector2.Vector{
		X: direction.X * velocity,
		Y: direction.Y * velocity,
	}

	rotationSpeed := -0.02 + rand.Float64()*0.04

	meteor := Meteor{
		position: pos,
		movement: movement,
		rotationSpeed: rotationSpeed,
		angle:    0,
		size:     size,
		sprite:   sprite,
		width:	  sprite.Bounds().Dx(),
		height:   sprite.Bounds().Dy(),
		done:	  false,	
	}
	// Add to list of meteors
	meteors = append(meteors, &meteor)

	return &meteor
}

// To create a new Meteor of size (0 to 3) at position x, y
func NewFragment(size int, x, y float64) *Meteor {
	if size > 3 {
		size = 3
	}
	sprite := meteorSprites[size]

	// select a random target	
	target := vector2.Vector{
		X: float64(rand.IntN(ScreenWidth)),
		Y: float64(rand.IntN(ScreenHeight)),
	}
	pos := vector2.Vector{
		X: float64(x),
		Y: float64(y),
	}

	// Randomized velocity
	velocity := 0.25 + rand.Float64()*1.5

	// Direction is the target minus the current position
	direction := vector2.Vector{
		X: (target.X - pos.X),
		Y: (target.Y - pos.Y),
	}

	// remove distance from Vector
	direction.Normalize()

	movement := vector2.Vector{
		X: direction.X * velocity,
		Y: direction.Y * velocity,
	}

	rotationSpeed := -0.02 + rand.Float64()*0.04

	meteor := Meteor{
		position: pos,
		movement: movement,
		rotationSpeed: rotationSpeed,
		angle:    0,
		size:     size,
		sprite:   sprite,
		width:	  sprite.Bounds().Dx(),
		height:   sprite.Bounds().Dy(),
		done:	  false,	
	}
	// Add to list of meteors
	meteors = append(meteors, &meteor)

	return &meteor
}

func (m *Meteor) ScreenPos() ScreenPos {
	return ScreenPos{
		X: 		m.position.X,
		Y: 		m.position.Y,
		Width: 	m.width,
		Height: m.height,
	}

}

func (m *Meteor) Update() {
	if !m.done {
		m.position.Add(m.movement)
		m.angle += m.rotationSpeed
		m.done =  m.position.X < -50 || m.position.X > ScreenWidth + 50 ||
				m.position.Y < -50 || m.position.Y > ScreenHeight + 50
	}
}

func (m *Meteor) Draw(screen *ebiten.Image) {
	if !m.done{
		DrawImage(screen, m.sprite, m.position.X, m.position.Y, m.angle)
	}
}

// Process meteor being hit and return score based on size
func (m *Meteor) Hit(explode bool) int{
	if !m.done {
		score := scores[m.size]
		m.size -= 1
		if m.size < 0 {
			m.done = true
			if explode {
				NewExplosion(m.position.X, m.position.Y, 20, 
					expColor, 0.075)
			}
		} else {
			m.sprite = meteorSprites[m.size]
			m.width = m.sprite.Bounds().Dx()
			m.height = m.sprite.Bounds().Dy()
			// split in 2 by creating a new fragment
			NewFragment(m.size, m.position.X, m.position.Y)
			// Create Explosion
			if explode {
				NewExplosion(m.position.X, m.position.Y, 20 * (m.size + 1), 
				expColor, 0.060 - (0.020 * float64(m.size)))
			}
			}
		return score
	} else {
		return 0
	}
}
