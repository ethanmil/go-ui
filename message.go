package ui

import (
	"image/color"
	"log"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Content struct {
	text       string
	suspension *time.Duration
	action     func()
}

const (
	messagePadding  = 12
	messageTextSize = 10
)

var (
	messageBGColor = color.White
	messageFont    = GetArcadeFont(messageTextSize)
)

type Message struct {
	content chan Content

	lines []string

	element *Element

	autoFormat bool
}

func NewMessage(msg string, action func() error, element *Element, autoFormat bool) *Message {
	var (
		width    = element.Width
		height   = element.Height
		textSize = text.BoundString(messageFont, msg)
	)

	// default message width to text size + padding unless specified in element
	if element.Width == 0 {
		width = textSize.Dx() + (messagePadding * 2)
		height = textSize.Dy() + (messagePadding * 2)

		// update element with new values
		element.Width = width
		element.Height = height
	}

	message := Message{
		content:    make(chan Content, 50),
		element:    element,
		autoFormat: autoFormat,
	}

	message.handleSetText()

	message.SetText(msg)

	return &message
}

func (m *Message) SetText(str string) {
	m.content <- Content{
		text: str,
	}
}

func (m *Message) Suspend(str string, t time.Duration) {
	m.content <- Content{
		text:       str,
		suspension: &t,
	}
}

func (m *Message) SuspendWithAction(str string, t time.Duration, a func()) {
	m.content <- Content{
		text:       str,
		suspension: &t,
		action:     a,
	}
}

func (m *Message) handleSetText() {
	charLimit := (m.element.Width - messagePadding) / messageTextSize

	go func() {
		for {
			content := <-m.content

			if content.suspension != nil {
				time.Sleep(*content.suspension)
			}

			currentLine := 0
			m.lines = []string{"", "", "", "", "", "", "", ""} // init eight lines
			if m.autoFormat && len(content.text) > 0 {
				content.text = autoFormat(content.text)
			}

			content.text += " " // extra, empty space for easier word parsing

			for i, char := range content.text {
				nextSpace := strings.Index(content.text[i:], " ")

				if len(m.lines[currentLine])+nextSpace >= charLimit {
					currentLine++
				}

				if len(m.lines) == currentLine {
					log.Fatal("current line should not excede line count")
				}

				m.lines[currentLine] += string(char)
				time.Sleep(25 * time.Millisecond)
			}

			if content.suspension != nil {
				time.Sleep(*content.suspension)
			}

			if content.action != nil {
				content.action()
			}
		}
	}()
}

func autoFormat(str string) string {
	runes := []rune(str)

	firstCharacter, lastCharacter := string(runes[0]), string(runes[len(runes)-1])

	runes[0] = []rune(strings.ToUpper(firstCharacter))[0]

	if lastCharacter != "." && lastCharacter != "!" && lastCharacter != "?" {
		runes = append(runes, '.')
	}

	return string(runes)
}

func (m *Message) Update() error {
	return nil
}

func (m *Message) Draw(img *ebiten.Image) {
	DrawRect(img, float64(m.element.X), float64(m.element.Y), float64(m.element.Width), float64(m.element.Height), messageBGColor)

	for i := range m.lines {
		text.Draw(img, m.lines[i], messageFont, m.element.X+messagePadding, 8+m.element.Y+messagePadding+((messageTextSize+12)*i), color.Black)
	}
}
