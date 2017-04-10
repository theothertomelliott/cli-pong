package main

import (
	"fmt"
	"time"

	termbox "github.com/nsf/termbox-go"
)

// Colors
const (
	backgroundColor = termbox.ColorBlack
	foregroundColor = termbox.ColorWhite
)

func main() {
	var closeEvent termbox.Event
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer func() {
		termbox.Close()
		if r := recover(); r != nil {
			fmt.Println("Recovered: ", r)
		} else {
			fmt.Println("Closed: ", closeEvent.Key)
		}
	}()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	width, height := termbox.Size()
	g := newGame(width, height)

	animationSpeed := 75 * time.Millisecond

	var paused bool

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch {
				case ev.Ch == 'p':
					paused = !paused
				case ev.Ch == 'w':
					g.moveUp(PlayerOneIndex)
				case ev.Ch == 's':
					g.moveDown(PlayerOneIndex)
				case ev.Key == termbox.KeyArrowUp:
					g.moveUp(PlayerTwoIndex)
				case ev.Key == termbox.KeyArrowDown:
					g.moveDown(PlayerTwoIndex)
				case ev.Key == termbox.KeyCtrlC:
					closeEvent = ev
					return
				}
			}
		default:
			if !paused {
				g.update()
			}
			render(g)
			time.Sleep(animationSpeed)
		}
	}
}

func render(g *game) {
	termbox.Clear(backgroundColor, backgroundColor)
	g.draw()
	termbox.Flush()
}
