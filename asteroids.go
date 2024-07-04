// Asteroids game written in Go using Ebitengine
// Author Paul Brace
// July 2024

package main

/* To do
Add music
*/

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 1000
	ScreenHeight = 800
	StartSpawnTime = 3			// New Meteor created every 3 seconds at start
	SpawnChangeTime = 0.10		// Gap between new meteors reduced by this every SpawnChangeInterval
	MinSpawnTime = 2			// Minimum gap between new meteors being created
	SpawnChangeInterval = 120   // spawn time changed every 2 minutes
)

// Game mode
const (
	Inst = 0
	InPlay = 1
	GameOver = 2
)

// Embeds all of asset resources to assets
//go:embed assets/*
var assets embed.FS

type Game struct{
	spawnSpeed			float64
	spawnTimer 			*Timer
	spawnUpdateTimer	*Timer
	playerHitTimer		*Timer
	player 				*Player
	scoreboard 			*ScoreBoard
	game_mode			int
}

func (g *Game) Update() error {
	UpdateStars()	// Background
	if g.game_mode == InPlay {
		UpdateAllTimers()
		ClearDoneExplosions()
		ClearDoneMeteors()
		ClearDoneMissiles()
		if g.spawnTimer.IsReady() {
			NewMeteor(g.player)
		}
		if g.spawnUpdateTimer.IsReady() {
			if g.spawnSpeed > MinSpawnTime {
				g.spawnSpeed -= SpawnChangeTime
				g.spawnTimer.ChangeTime(g.spawnSpeed, true)
			}
		}
		if g.player.alive {
			g.player.Update()
		}
		UpdateAllMeteors()
		UpdateAllMissiles()
		for _, miss := range missiles {
			// Check if hit a meteor
			missPos := miss.ScreenPos()
			for _, met := range meteors {
				if missPos.Collides(met.ScreenPos()){
					// mark as hit, update score and set as done so removed next frame
					g.scoreboard.score += met.Hit(true)
					miss.done = true
					break
				}
			}
		}
		if g.player.alive {
			playerPos := g.player.ScreenPos()
			for _, met := range(meteors){
				// Check if hit player
				if playerPos.Collides(met.ScreenPos()){
					g.player.Hit()
					g.scoreboard.lives -= 1
					g.playerHitTimer = NewTimer(3, false)
					g.spawnTimer.active = false
					met.Hit(false)
					break
				}
			}
		} else {
			if g.playerHitTimer.IsReady() {
				if g.scoreboard.lives == 0 {
					g.game_mode = GameOver
					g.spawnTimer.active = false
					g.spawnUpdateTimer.active = false
				} else {
					g.player.Reset()
					ClearAllMeteors()
					ClearAllMissiles()
					ClearAllExplosions()
					g.spawnTimer.Reset()
 		   			// create a single meteor to start
					NewMeteor(g.player)
				}
			}
		}
		UpdateAllExplosions()
	} else {
		// Check if a key has been pressed
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			// Stop player firing for reload period ans set for new game
			g.player.Reset()
			g.player.loaded = false
			reloadTimer = reloadTime
			ClearAllMeteors()
			ClearAllMissiles()
			ClearAllExplosions()
			g.scoreboard.score = 0
			g.scoreboard.lives = 3
			g.spawnSpeed = StartSpawnTime
			g.spawnTimer.ChangeTime(StartSpawnTime, true)
			g.spawnUpdateTimer.Reset()
			g.scoreboard.LoadHighScore()
			g.scoreboard.highScoreSaved = false
   			// create a single meteor to start
			NewMeteor(g.player)
			g.game_mode = InPlay
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	DrawStars(screen)
	switch g.game_mode {
	case InPlay:
		g.player.Draw(screen)
		DrawAllMeteors(screen)
		DrawAllMissiles(screen)
		DrawAllExplosions(screen)
		g.scoreboard.DrawScore(screen)
	case GameOver:
		g.scoreboard.DrawGameOver(screen)
	default:
		g.scoreboard.DrawInstructions(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	g := &Game{
		spawnSpeed: StartSpawnTime,
		spawnTimer: NewTimer(StartSpawnTime, true),
		spawnUpdateTimer: NewTimer(SpawnChangeInterval, true),
		player: NewPlayer(),
		scoreboard: NewScoreBoard(),
		game_mode: Inst,
	}
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Asteroids")
	CreateStarField()
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}