package ui

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// interface check
var _ Sprite = &Image{}

type Image struct {
	img *ebiten.Image

	element *Element
}

func NewImage(path string, element *Element) *Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return &Image{img, element}
}

func NewImageFromImage(img *Image, element *Element) *Image {
	return &Image{img.img, element}
}

func (i *Image) Update() error {
	return nil
}

func (i *Image) Draw(img *ebiten.Image) {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(i.element.X), float64(i.element.Y))

	img.DrawImage(i.img, &op)
}
