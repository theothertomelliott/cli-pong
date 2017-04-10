package main

import termbox "github.com/nsf/termbox-go"

type canvas struct {
	width, height    int
	offsetX, offsetY int
	backgroundColor  termbox.Attribute
	foregroundColor  termbox.Attribute
}

type TextAlignment int

const (
	TextAlignmentLeft   TextAlignment = iota
	TextAlignmentRight  TextAlignment = iota
	TextAlignmentCenter TextAlignment = iota
)

func newCanvas(x, y, width, height int, foreground, background termbox.Attribute) canvas {
	return canvas{
		offsetX:         x,
		offsetY:         y,
		width:           width,
		height:          height,
		foregroundColor: foreground,
		backgroundColor: background,
	}
}

func (c *canvas) drawText(x, y int, text string, alignment TextAlignment) {
	var minX int
	switch alignment {
	case TextAlignmentCenter:
		minX = x - (len(text) / 2)
	case TextAlignmentRight:
		minX = x - len(text)
	default:
		minX = x
	}
	for i, r := range text {
		c.setCell(minX+i, y, r)
	}
}

func (c *canvas) setCell(x, y int, r rune) {
	if x < 0 || x > c.width {
		return
	}
	if y < 0 || y >= c.height {
		return
	}
	termbox.SetCell(c.offsetX+x, c.offsetY+y, r, c.foregroundColor, c.backgroundColor)
}

func (c *canvas) clear() {
	for x := 0; x < c.width; x++ {
		for y := 0; y < c.height; y++ {
			termbox.SetCell(c.offsetX+x, c.offsetY+y, ' ', c.foregroundColor, c.backgroundColor)
		}
	}
}
