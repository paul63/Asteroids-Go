// Sprite/image functions 
// Screen position struct and methods
// for games written in Go using Ebitengine
// Author Paul Brace
// July 2024

package main

import (
	"image"
	_ "image/png"
	"math"
	"github.com/hajimehoshi/ebiten/v2"
)

// Load the image requested in name
func LoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

// Draws sprite on the screen center at x, y and rotated by rotate
func DrawImage(screen  *ebiten.Image, sprite *ebiten.Image, x, y, rotate float64) {
	op := &ebiten.DrawImageOptions{}
	width := sprite.Bounds().Dx()
	height := sprite.Bounds().Dy()
	halfW := float64(width / 2)
	halfH := float64(height / 2)
	// move image so center aligns with 0, 0
	op.GeoM.Translate(-halfW, -halfH)
	// do the rotation
	op.GeoM.Rotate(rotate)
	// move it to required position X & Y will be center of sprite as relative to 0,0
	op.GeoM.Translate(x , y)
	screen.DrawImage(sprite, op)

}

// Struct representing the sprite position on the screen
// x and y are center of image
type  ScreenPos struct {
	X      float64
	Y      float64
	Width  int
	Height int
}

// Create a new screen position using parameters passed
func NewScreenPos(x, y float64, width, height int) ScreenPos {
	return ScreenPos{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// Checks if one screen position intersects another
func (my ScreenPos) Intersects(other ScreenPos) bool {
	return 	(math.Abs(my.X - other.X) < float64(my.Width) / 2 + float64(other.Width) / 2) && 
			(math.Abs(my.Y - other.Y) < float64(my.Height) / 2 + float64(other.Height) / 2)
}

// Check if distance between 2 screen positions is less than sum of radiuses (radius = half with)
// 2 pixel margin deducted so slight overlap
func (my ScreenPos) Collides(other ScreenPos) bool {
	distance := math.Sqrt((other.X - my.X) * (other.X - my.X) + (other.Y - my.Y) * (other.Y - my.Y))
	return distance < float64(my.Width) / 2 + float64(other.Width) / 2 - 2	
}
