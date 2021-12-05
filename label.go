package ui

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

const (
	labelPadding = 2
	labelBorder  = 1
)

type Label struct {
	text string

	element *Element

	fontSize int
	fontFace font.Face
	color    color.Color

	animationLabel *Label

	show            bool
	cancelAnimation chan int
}

type LabelOption struct {
	FontSize int
	Color    color.Color
}

func NewLabel(msg string, element *Element, op *LabelOption) *Label {
	fontSize := defaultFontSize
	fontFace := defaultFont
	col := color.Color(defaultColor)

	if op != nil {
		if op.FontSize != 0 {
			fontSize = op.FontSize
			fontFace = GetArcadeFont(fontSize)
		}
		if op.Color != nil {
			col = op.Color
		}
	}

	var (
		width    = element.Width
		height   = element.Height
		textSize = text.BoundString(fontFace, msg)
	)

	// default label width to text size + padding unless specified in element
	if element.Width == 0 {
		width = textSize.Dx() + (labelPadding * 2)
		height = textSize.Dy() + (labelPadding * 2)

		// update element with new values
		element.Width = width
		element.Height = height
	}

	return &Label{msg, element, fontSize, fontFace, col, nil, true, make(chan int)}
}

func (l *Label) SetText(str string) {
	l.text = str
}

func (l *Label) SetTextWithAnimation(newString string, animationString string, col color.Color, dir Direction) {
	l.SetText(newString)

	if l.animationLabel != nil && l.animationLabel.show {
		l.cancelAnimation <- 0
	}

	l.animationLabel = NewLabel(
		animationString,
		&Element{X: l.element.X + (len(newString) * l.fontSize) - (len(animationString) * l.fontSize), Y: l.element.Y},
		&LabelOption{FontSize: l.fontSize, Color: col},
	)

	frameRate := time.Millisecond * 30
	frameTimer := time.NewTimer(frameRate)

	movementDelta := 30
	movementIterator := 1

	go func() {
		for movementIterator < movementDelta {
			select {
			case <-l.cancelAnimation:
				l.animationLabel.show = false
				return
			case <-frameTimer.C:
				l.animationLabel.show = true

				movementIterator++

				if movementIterator == movementDelta {
					l.animationLabel.show = false
					return
				}

				var x, y int
				switch dir {
				case Up:
					y = -1
				case Down:
					y = 1
				case Right:
					x = 1
				case Left:
					x = -1
				}

				l.animationLabel.element.X += x
				l.animationLabel.element.Y += y

				r, g, b, _ := l.animationLabel.color.RGBA()
				a := uint32(255 / movementDelta * (movementDelta - movementIterator))
				l.animationLabel.color = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}

				frameTimer.Reset(frameRate)
			}
		}
	}()
}

func (l *Label) Update() error {
	return nil
}

func (l *Label) Draw(img *ebiten.Image) {
	if !l.show {
		return
	}

	labelTextSize := text.BoundString(l.fontFace, l.text)
	text.Draw(img, l.text, l.fontFace, l.element.X+labelPadding, labelTextSize.Dy()+l.element.Y+labelPadding, l.color)

	if l.animationLabel != nil {
		animationLabelTextSize := text.BoundString(l.animationLabel.fontFace, l.animationLabel.text)
		text.Draw(img, l.animationLabel.text, l.animationLabel.fontFace, l.animationLabel.element.X+labelPadding, animationLabelTextSize.Dy()+l.animationLabel.element.Y+labelPadding, l.animationLabel.color)
	}
}
