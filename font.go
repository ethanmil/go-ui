package ui

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	defaultFontSize = 12
	defaultFont     = GetArcadeFont(defaultFontSize)
	defaultColor    = color.White
)

func GetArcadeFont(size int) font.Face {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}

	arcadeFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	return arcadeFont
}
