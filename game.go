package main

import (
	"fmt"
	"math"
	"time"

	termbox "github.com/nsf/termbox-go"
)

const (
	PlayerOneIndex = 0
	PlayerTwoIndex = 1

	PlayerOneRune = '▒'
	PlayerTwoRune = '▓'
	BallRune      = '█'
)

type vector struct {
	x, y int
}

type sprite struct {
	pos   vector
	speed vector
}

type ball struct {
	sprite
}

type paddle struct {
	sprite
	size int
	r    rune
}

type game struct {
	ball        ball
	scores      []int
	players     []*paddle
	paddleSpeed int
	gameCanvas  canvas
	scoreCanvas canvas
	hintCanvas  canvas
}

func newGame(width, height int) *game {
	g := game{}
	g.scoreCanvas = newCanvas(0, 0, width, 1, termbox.ColorWhite, termbox.ColorBlack)
	g.gameCanvas = newCanvas(0, 1, width, height-2, termbox.ColorWhite, termbox.ColorBlue)
	g.hintCanvas = newCanvas(0, height-1, width, 1, termbox.ColorWhite, termbox.ColorBlack)
	g.ball.speed.x, g.ball.speed.y = 1, 1
	g.centerBall()
	g.resetBall()

	g.players = make([]*paddle, 2)
	g.scores = make([]int, 2)

	playerSize := 4
	g.paddleSpeed = 2
	playerStartingY := (height / 2) - (playerSize / 2)

	g.players[PlayerOneIndex] = &paddle{
		sprite: sprite{
			pos: vector{x: 1, y: playerStartingY},
		},
		size: playerSize,
		r:    PlayerOneRune,
	}
	g.players[PlayerTwoIndex] = &paddle{
		sprite: sprite{
			pos: vector{x: width - 2, y: playerStartingY},
		},
		size: playerSize,
		r:    PlayerTwoRune,
	}

	return &g
}

func (g *game) moveUp(playerIndex int) {
	player := g.players[playerIndex]
	if player.pos.y > 0 {
		player.pos.y -= g.paddleSpeed
	}
	if player.pos.y <= 0 {
		player.pos.y = 0
	}
}

func (g *game) moveDown(playerIndex int) {
	player := g.players[playerIndex]
	if (player.pos.y + player.size) < g.gameCanvas.height {
		player.pos.y += g.paddleSpeed
	}
	if (player.pos.y + player.size) > g.gameCanvas.height {
		player.pos.y = g.gameCanvas.height - player.size
	}
}

func (g *game) resetBall() {
	g.centerBall()
	newVelocity := g.ball.speed
	newVelocity.x, newVelocity.y = -newVelocity.x, -newVelocity.y
	g.ball.speed.x, g.ball.speed.y = 0, 0

	go func() {
		time.Sleep(1 * time.Second)
		g.ball.speed = newVelocity
	}()
}

func (g *game) centerBall() {
	g.ball.pos.x = g.gameCanvas.width / 2
	g.ball.pos.y = g.gameCanvas.height / 2
}

func (g *game) update() {
	g.ball.pos.x += g.ball.speed.x
	g.ball.pos.y += g.ball.speed.y

	if g.ball.pos.x >= g.gameCanvas.width {
		g.scores[PlayerOneIndex]++
		g.resetBall()
	}
	if g.ball.pos.x <= 0 {
		g.scores[PlayerTwoIndex]++
		g.resetBall()
	}
	if g.ball.pos.y >= g.gameCanvas.height-1 || g.ball.pos.y <= 0 {
		g.ball.speed.y = -g.ball.speed.y
	}

	for _, player := range g.players {
		if math.Abs(float64(g.ball.pos.x-player.pos.x)) == 1 {
			if g.ball.pos.y >= player.pos.y &&
				g.ball.pos.y <= player.pos.y+player.size {
				g.ball.speed.x = -g.ball.speed.x
			}
		}
	}

}

func (g *game) draw() {
	g.gameCanvas.clear()

	g.drawBall()
	g.drawPaddles()
	g.drawScores()
	g.drawHints()
}

func (g *game) drawScores() {
	g.scoreCanvas.drawText(g.gameCanvas.width/2, 0, fmt.Sprintf("%v : %v", g.scores[PlayerOneIndex], g.scores[PlayerTwoIndex]), TextAlignmentCenter)
}

func (g *game) drawHints() {
	g.hintCanvas.drawText(0, 0, "Up: W, Down: S", TextAlignmentLeft)
	g.hintCanvas.drawText(g.hintCanvas.width/2, 0, "Pause: P   Quit: Ctrl+C", TextAlignmentCenter)
	g.hintCanvas.drawText(g.hintCanvas.width, 0, "Up: ⬆️, Down: ⬇️", TextAlignmentRight)
}

func (g *game) drawBall() {
	g.gameCanvas.setCell(g.ball.pos.x, g.ball.pos.y, BallRune)
}

func (g *game) drawPaddles() {
	for _, player := range g.players {
		for i := 0; i < player.size; i++ {
			g.gameCanvas.setCell(player.pos.x, player.pos.y+i, player.r)
		}
	}

}
