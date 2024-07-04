// Explosion struct and methods
// for games written in Go using Ebitengine
// Author Paul Brace
// July 2024

package main

import (
	"image/color"
	"math/rand/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/paul63/vector2"
)

var explosions [] *Explosion

func UpdateAllExplosions(){
	for _, e := range explosions {
		e.Update()
	} 
}

func DrawAllExplosions(screen *ebiten.Image){
	for _, e := range explosions {
		e.Draw(screen)
	} 
}

func ClearDoneExplosions(){
	for i, e := range(explosions){
		if e.done{
			explosions = append(explosions[:i], explosions[i+1:]...)
			break
		}		
	}
}

func ClearAllExplosions(){
	explosions = nil
}

type Explosion struct{
	x		float64
	y   	float64
	size 	int
	color   color.Color
	rate 	float64
	done 	bool
	particles [] *Particle
}

// Create a new explosion and add to list
// x, y = Center of explosion. int the number of particles. 
// Color = color of particles. rate = the rate that the prticles fade away
func NewExplosion( x_pos, y_pos float64, size int, color color.Color, rate float64) *Explosion{
	var particles [] *Particle
	// create particles of random size and direction
	for i := 0; i < size * 2; i++{
		part := NewParticle(x_pos, y_pos, 
							rand.Float64() * 4,
							color,
							(rand.Float64() - 0.5) * rand.Float64() * 6,
							(rand.Float64() - 0.5) * rand.Float64() * 6,
							rate)
		particles = append(particles, part)
	}
	exp := Explosion{
			x: x_pos,
			y: y_pos,
			size: size,
			color: color,
			rate: rate,
			done: false,
			particles: particles,
	}
	explosions = append(explosions, &exp)
	return &exp
}

func (e *Explosion) Update(){
	for i, p := range(e.particles){
		p.Update()
		if p.done{
			e.particles = append(e.particles[:i], e.particles[i+1:]...)
			break
		}		
	}
	if len(e.particles) == 0 {
		e.done = true
	}
}

func (e *Explosion) Draw(screen *ebiten.Image){
	for _, p := range(e.particles){
		if !p.done {
			p.Draw(screen)
		}
	}

}

const Friction = 0.99  // The rate that the particles slow down

type Particle struct{
	position  	vector2.Vector
	radius 		float64
	color 		color.Color
	velocity 	vector2.Vector
	rate 		float64
	done 		bool
}

func NewParticle(x, y, radius float64, color color.Color, velocityX, velocityY, rate float64) *Particle{
	part := Particle{
		position: vector2.Vector{X:x, Y:y,},
		radius: radius,
		color: color,
		velocity: vector2.Vector{X:velocityX, Y:velocityY},
		rate: rate,
		done: false,
	}
	return &part
}

func (p *Particle) Update(){
	p.velocity.X *= Friction
	p.velocity.Y *= Friction
	if p.radius > 0{
		p.radius -= p.rate
	}
	if p.radius <= 0 {
		p.done = true
	}
	if !p.done{
		p.position.Add(p.velocity)
	}
}

func (p *Particle) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, float32(p.position.X), float32(p.position.Y), float32(p.radius), p.color, false)
}