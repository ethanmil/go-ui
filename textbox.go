package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

const (
	textboxPadding   = 12
	fontPaddingRatio = 0.6
)

type Textbox struct {
	runes []rune
	msg   string

	keyStrokes chan Key

	element *Element

	font font.Face
}

func NewTextbox(msg string, keyStrokes chan Key, element *Element) *Textbox {
	tb := &Textbox{
		runes:      make([]rune, 100),
		msg:        msg,
		keyStrokes: keyStrokes,
		element:    element,
		font:       GetArcadeFont(element.Height - int(float64(element.Height)*fontPaddingRatio)),
	}

	tb.handleKeyStrokes()

	return tb
}

func (t *Textbox) handleKeyStrokes() {
	go func() {
		for {
			key := <-t.keyStrokes

			if key.key == ebiten.KeyBackspace {
				if len(t.msg) >= 1 {
					t.msg = t.msg[:len(t.msg)-1]
				}
				continue
			}

			t.msg += key.display
		}
	}()
}

func (t *Textbox) GetText() string {
	return t.msg
}

func (t *Textbox) Update() error {
	t.runes = ebiten.AppendInputChars(t.runes[:0])
	t.msg += string(t.runes)

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if len(t.msg) >= 1 {
			t.msg = t.msg[:len(t.msg)-1]
		}
	}

	return nil
}

func (t *Textbox) Draw(img *ebiten.Image) {
	DrawRect(img, float64(t.element.X), float64(t.element.Y), float64(t.element.Width), float64(t.element.Height), color.White)
	text.Draw(img, t.msg, t.font, t.element.X+textboxPadding, t.element.Y+t.element.Height-(int(float64(t.element.Height)*(fontPaddingRatio/2))), color.Black)
}
