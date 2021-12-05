package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var letterKeySet = []ebiten.Key{
	ebiten.KeyQ, ebiten.KeyW, ebiten.KeyE, ebiten.KeyR, ebiten.KeyT, ebiten.KeyY, ebiten.KeyU, ebiten.KeyI, ebiten.KeyO, ebiten.KeyP,
	ebiten.KeyA, ebiten.KeyS, ebiten.KeyD, ebiten.KeyF, ebiten.KeyG, ebiten.KeyH, ebiten.KeyJ, ebiten.KeyK, ebiten.KeyL,
	ebiten.KeyZ, ebiten.KeyX, ebiten.KeyC, ebiten.KeyV, ebiten.KeyB, ebiten.KeyN, ebiten.KeyM, ebiten.KeyBackspace,
}

var numberKeySet = []ebiten.Key{
	ebiten.KeyDigit7, ebiten.Key8, ebiten.Key9,
	ebiten.Key4, ebiten.Key5, ebiten.Key6,
	ebiten.Key1, ebiten.Key2, ebiten.Key3,
	ebiten.Key0, ebiten.KeyBackspace,
}

type Keyboard struct {
	element *Element

	keys []*Key

	keyStrokes chan Key
}

func NewLetterKeyboard(element *Element) *Keyboard {
	keys := []*Key{}

	keyStrokes := make(chan Key, 10)

	keyWidth := element.Width / 10
	keyHeight := element.Height / 3

	for i, key := range letterKeySet {
		extraWidth := 0
		extraSpace := 0
		if key == ebiten.KeyBackspace {
			extraWidth = keyWidth * 2
			extraSpace = keyWidth
		}

		keys = append(keys, NewKey(
			key,
			key.String(),
			keyStrokes,
			&Element{
				X:      element.X + (int(i%10) * keyWidth) + extraSpace,
				Y:      element.Y + (int(i/10) * keyHeight),
				Width:  keyWidth + extraWidth,
				Height: keyHeight,
			},
		))
	}

	return &Keyboard{
		element:    element,
		keyStrokes: keyStrokes,
		keys:       keys,
	}
}

func NewNumberKeyboard(element *Element) *Keyboard {
	keys := []*Key{}

	keyStrokes := make(chan Key, 10)

	keyWidth := element.Width / 3
	keyHeight := element.Height / 4

	for i, key := range numberKeySet {
		extraWidth := 0
		extraSpace := 0
		displayString := key.String()[5:6]
		if key == ebiten.KeyBackspace {
			extraWidth = keyWidth
			displayString = "Backspace"
		}

		keys = append(keys, NewKey(
			key,
			displayString,
			keyStrokes,
			&Element{
				X:      element.X + (int(i%3) * keyWidth) + extraSpace,
				Y:      element.Y + (int(i/3) * keyHeight),
				Width:  keyWidth + extraWidth,
				Height: keyHeight,
			},
		))
	}

	return &Keyboard{
		element:    element,
		keyStrokes: keyStrokes,
		keys:       keys,
	}
}

func (k *Keyboard) GetKeyStrokes() chan Key {
	return k.keyStrokes
}

func (k *Keyboard) Update() error {
	for _, key := range k.keys {
		key.Update()
	}

	return nil
}

func (k *Keyboard) Draw(img *ebiten.Image) {
	for _, key := range k.keys {
		key.Draw(img)
	}
}

type Key struct {
	key ebiten.Key

	display string

	keyStrokes chan Key

	element *Element
}

func NewKey(key ebiten.Key, display string, keyStrokes chan Key, element *Element) *Key {
	return &Key{key, display, keyStrokes, element}
}

func (k *Key) Update() error {
	mx, my := 0, 0

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my = ebiten.CursorPosition()
	}

	touchIDs := []ebiten.TouchID{}
	touchIDs = inpututil.AppendJustPressedTouchIDs(touchIDs)
	if len(touchIDs) > 0 {
		mx, my = ebiten.TouchPosition(touchIDs[0])
	}

	if k != nil && mx >= k.element.X && mx <= k.element.X+k.element.Width &&
		my >= k.element.Y && my <= k.element.Y+k.element.Height {
		k.keyStrokes <- *k
	}

	return nil
}

func (k *Key) Draw(img *ebiten.Image) {
	textSize := text.BoundString(defaultFont, k.display)
	DrawRect(img, float64(k.element.X+2), float64(k.element.Y+2), float64(k.element.Width-4), float64(k.element.Height-4), color.White)
	text.Draw(img, k.display, defaultFont, k.element.X+(k.element.Width/2)-(textSize.Dx()/2), k.element.Y+(k.element.Height/2)-(textSize.Dy()/2)+textSize.Dy(), color.Black)
}
