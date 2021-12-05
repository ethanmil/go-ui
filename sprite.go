package ui

import "github.com/hajimehoshi/ebiten/v2"

type Sprite interface {
	Update() error
	Draw(img *ebiten.Image)
}
