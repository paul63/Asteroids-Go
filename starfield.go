// Star background functions, struct and methods
// for Asteroids written in Go using Ebitengine
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

const (
	NumStars = 50
	StarSpeed = 0.25
)

var stars [] *Star

var star_white color.Color = color.RGBA{125, 125, 125, 75}

func CreateStarField() {
	if len(stars) != 0 {
		// clear existing stars
		stars = nil
	}
	for i := 0; i < NumStars; i++ {
		s := NewStar(
			float64(rand.IntN(ScreenWidth)),
			float64(rand.IntN(ScreenHeight)),
			rand.IntN(3) + 1)
		stars = append(stars, s)
	}
}

func UpdateStars(){
	for _, s := range(stars) {
		s.Update()
	}
}

func DrawStars(screen *ebiten.Image){
	for _, s := range(stars) {
		s.Draw(screen)
	}
}

type Star struct {
	position	vector2.Vector
	velocity 	vector2.Vector
	radius		int
}

func NewStar(x, y float64, size int) *Star{
	s := Star{
		position: vector2.Vector{X: x, Y: y},
		velocity: vector2.Vector{X: 0, Y:StarSpeed},
		radius: size,
	}
	stars = append(stars, &s)
	return &s
}


func (s *Star) Draw(screen *ebiten.Image){
	vector.DrawFilledCircle(screen, float32(s.position.X), float32(s.position.Y), float32(s.radius), star_white, false)
}

func (s *Star) Update() {
	s.position.Add(s.velocity)
	if s.position.Y > ScreenHeight {
		s.position.Y = 0
		s.position.X = float64(rand.IntN(ScreenWidth))
	}
}

