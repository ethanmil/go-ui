package ui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	emptyImage    = ebiten.NewImage(3, 3)
	emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	emptyImage.Fill(color.White)
}

func DrawRect(dst *ebiten.Image, x, y, width, height float64, clr color.Color) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(width, height)
	op.GeoM.Translate(x, y)
	op.ColorM.Scale(colorToScale(clr))
	dst.DrawImage(emptySubImage, op)
}

func colorToScale(clr color.Color) (float64, float64, float64, float64) {
	cr, cg, cb, ca := clr.RGBA()
	if ca == 0 {
		return 0, 0, 0, 0
	}
	return float64(cr) / float64(ca), float64(cg) / float64(ca), float64(cb) / float64(ca), float64(ca) / 0xffff
}
