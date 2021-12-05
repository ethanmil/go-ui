package ui

import (
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Direction string

const (
	Up    Direction = "up"
	Down  Direction = "down"
	Right Direction = "right"
	Left  Direction = "left"
)

// interface check
var _ Sprite = &Animation{}

type Animation struct {
	image *ebiten.Image

	frameTotal  int
	frameWidth  int
	frameHeight int

	frameRate    time.Duration
	cycleTotal   int
	executeAtEnd func()

	cycleCount    int
	currentFrame  int
	showNextFrame bool

	element *Element
}

func NewAnimation(image *ebiten.Image, frameTotal, frameWidth, frameHeight int, frameRate time.Duration, cycleTotal int, element *Element, executeAtEnd func()) *Animation {
	animation := Animation{
		image:        image,
		frameTotal:   frameTotal,
		frameWidth:   frameWidth,
		frameHeight:  frameHeight,
		frameRate:    frameRate,
		cycleTotal:   cycleTotal,
		executeAtEnd: executeAtEnd,
		element:      element,
	}

	frameTimer := time.NewTimer(frameRate)
	go func() {
		for {
			<-frameTimer.C
			animation.showNextFrame = true
			frameTimer.Reset(frameRate)
		}
	}()

	return &animation
}

func (a *Animation) Update() error {
	if a.showNextFrame {
		a.currentFrame++
		if a.currentFrame >= a.frameTotal {
			a.currentFrame = 0
			a.cycleCount++
		}

		a.showNextFrame = false
	}

	if a.cycleCount == a.cycleTotal {
		a.executeAtEnd()
	}

	return nil
}

func (a *Animation) Draw(img *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(a.element.X), float64(a.element.Y))
	sx, sy := a.currentFrame*a.frameWidth, 0
	img.DrawImage(a.image.SubImage(image.Rect(sx, sy, sx+a.frameWidth, sy+a.frameHeight)).(*ebiten.Image), op)
}
