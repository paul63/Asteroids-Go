// Scoreboard struct and methods
// for Asteroids written in Go using Ebitengine
// Author Paul Brace
// July 2024

package main

import (
	"bytes"
	"fmt"
	"image/color"
	"os"
	_ "embed"
	"strconv"	
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Embeds file specified into scoreRegular
//go:embed fonts/kenney-future.ttf
var scoreFont []byte

// Embeds instruction font
//go:embed fonts/FiraSans-Regular.ttf
var instFont []byte

var (
	scoreFace *text.GoTextFaceSource
	instFace *text.GoTextFaceSource
	green color.Color = color.RGBA{0, 255, 0, 255}
	yellow color.Color =  color.RGBA{255, 255, 0,255}
	white color.Color =  color.RGBA{255, 255, 255, 255}
	aqua color.Color =  color.RGBA{0, 255, 255, 255}
)

func LoadFont(fontStream []byte) *text.GoTextFaceSource{
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fontStream))
	if err != nil {
		panic(err)
	}
	return s
}

type ScoreBoard struct {
	score int
	highScore int
	lives	int
	highScoreSaved	bool
}

func (sb *ScoreBoard) LoadHighScore() int {
	// Load high score - ignore errors - if err then HS set to 0
    buff, err := os.ReadFile("score.txt")	
	hs := 0
	if err == nil{
		hs, _ = strconv.Atoi(string(buff))
	}
	sb.highScore = hs
	return hs
}

func NewScoreBoard() *ScoreBoard{
	sb := ScoreBoard {
		score: 0,
		highScore: 0,
		lives: 3,
		highScoreSaved: false,
	}
	sb.LoadHighScore()
	// Load fonts
	scoreFace = LoadFont(scoreFont)
	instFace = LoadFont(instFont)
	return &sb
}

func (sb *ScoreBoard) SaveHighScore(){
	if !sb.highScoreSaved {
		// Save new high score
		score := []byte(fmt.Sprint(sb.score))
		// Note 0644 is there for file create
		// 	  The file's owner can read and write (6)
		//    Users in the same group as the file's owner can read (first 4)
		//    All users can read (second 4)
		err := os.WriteFile("score.txt", score, 0644)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Unable to write to file.")
		}
		sb.highScoreSaved = true
	}
}

func (sb *ScoreBoard) DrawScore(screen *ebiten.Image){
	op := &text.DrawOptions{}
	op.GeoM.Translate(20, 20)
	text.Draw(screen, fmt.Sprintf("Score: %06d", sb.score), &text.GoTextFace{
		Source: scoreFace,
		Size:   20,
	}, op)	
	op = &text.DrawOptions{}
	op.GeoM.Translate(300, 20)
	text.Draw(screen, fmt.Sprintf("High Score: %06d",sb.highScore), &text.GoTextFace{
		Source: scoreFace,
		Size:   20,
	}, op)
	for i := 0; i < sb.lives; i++ {
		DrawImage(screen, playerSprite, float64(ScreenWidth - 35 - i * 40), 30, 0)
	}	
}

func (sb *ScoreBoard) DrawCenter(screen *ebiten.Image, s string, x, y, size int, color color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.PrimaryAlign = text.AlignCenter
	op.ColorScale.ScaleWithColor(color)
	text.Draw(screen, s, &text.GoTextFace{
		Source: instFace,
		Size:   float64(size),
	}, op)	
}

func (sb *ScoreBoard) DrawLeft(screen *ebiten.Image, s string, x, y, size int, color color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.LineSpacing = 30
	op.ColorScale.ScaleWithColor(color)
	text.Draw(screen, s, &text.GoTextFace{
		Source: instFace,
		Size:   float64(size),
	}, op)	
}

func (sb *ScoreBoard) DrawInstructions(screen *ebiten.Image){

	instructions := `Destroy the asteroids before they hit you.
A new asteroid appears every 3 seconds but the frequency increases
as the game progresses.
You have 3 lives.

Easy mode:
    Position the mouse pointer in the direction of the asteroid
    and press the left mouse button to fire.
    Right button to move.
    Middle button to hyperjump.

Hard mode. Press:
    Left and right arrow to rotate ship.
    Down arrow to reverse direction of ship.
    Up arrow to move.
    H to hyperjump.
    Space to fire.

You can have multiple missiles flying at one time.
Score for hitting an asteroid:
    Large = 25 Medium = 50 Small = 75 Tiny = 100 points`

	sb.DrawCenter(screen, "Asteroids", ScreenWidth/2, 20, 40, yellow)	
	sb.DrawLeft(screen, instructions, 200, 100, 20, white)
	sb.DrawCenter(screen, "Press space bar to play", ScreenWidth/2, 735, 30, aqua)
}

func (sb *ScoreBoard) DrawGameOver(screen *ebiten.Image){
	sb.DrawCenter(screen, "Asteroids", ScreenWidth/2, 20, 40, yellow)	
	sb.DrawCenter(screen, "Game Over", ScreenWidth/2, 200, 40, white)	
	sb.DrawCenter(screen, fmt.Sprintf("Your Score: %06d", sb.score), ScreenWidth/2, 400, 40, white)	
	if sb.score > sb.highScore {
		sb.DrawCenter(screen, "Congratulations a new high score", ScreenWidth/2, 500, 60, green)
		sb.SaveHighScore()	
	}
	sb.DrawCenter(screen, "Press space bar to play again", ScreenWidth/2, 700, 30, aqua)
}
