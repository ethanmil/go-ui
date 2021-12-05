package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

const (
	buttonPadding      = 6
	buttonPaddingRatio = 0.5
)

type Button struct {
	Text   string
	Action func(b *Button) error

	Element *Element

	labelElement *Element

	font font.Face
}

func NewButton(msg string, action func(b *Button) error, element *Element) *Button {
	var (
		width        = element.Width
		height       = element.Height
		font         = GetArcadeFont(element.Height - int(float64(element.Height)*fontPaddingRatio))
		textSize     = text.BoundString(font, msg)
		labelElement = &Element{}
	)

	// default button width to text size + padding unless specified in element
	if element.Width == 0 {
		width = textSize.Dx() + (buttonPadding * 2)
		height = textSize.Dy() + (buttonPadding * 2)

		// update element with new values
		element.Width = width
		element.Height = height
	}

	labelElement.X = element.X + (element.Width / 2) - (textSize.Dx() / 2)
	labelElement.Y = element.Y + (element.Height / 2) + (textSize.Dy() / 2)

	return &Button{msg, action, element, labelElement, font}
}

func (b *Button) SetText(str string) {
	b.Text = str
}

func (b *Button) Update() error {
	mx, my := 0, 0

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my = ebiten.CursorPosition()
	}

	touchIDs := []ebiten.TouchID{}
	touchIDs = inpututil.AppendJustPressedTouchIDs(touchIDs)
	if len(touchIDs) > 0 {
		mx, my = ebiten.TouchPosition(touchIDs[0])
	}

	if b != nil && mx >= b.Element.X && mx <= b.Element.X+b.Element.Width &&
		my >= b.Element.Y && my <= b.Element.Y+b.Element.Height {
		return b.Action(b)
	}

	return nil
}

func (b *Button) Draw(img *ebiten.Image) {
	DrawRect(img, float64(b.Element.X), float64(b.Element.Y), float64(b.Element.Width), float64(b.Element.Height), color.White)
	text.Draw(img, b.Text, b.font, b.labelElement.X, b.labelElement.Y, color.Black)
}
