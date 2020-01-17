package kitchen

import (
	"image/color"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

type (
	// DebugText ...
	DebugText struct {
		text              string
		uiFont            font.Face
		uiFontColor       color.Color
		uiFontColorLooper *Looper
		uiFontMHeight     int
		x                 int
		y                 int
	}
)

// NewDebugText ...
func NewDebugText() (*DebugText, error) {
	// ebitenフォントのテスト
	tt, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	uiFont := truetype.NewFace(tt, &truetype.Options{
		Size:    48,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	b, _, _ := uiFont.GlyphBounds('M')
	uiFontMHeight := (b.Max.Y - b.Min.Y).Ceil()

	s := &DebugText{
		uiFont:            uiFont,
		uiFontColor:       color.White,
		uiFontColorLooper: NewLooper(255, 1, 200, 255),
		uiFontMHeight:     uiFontMHeight,
		x:                 50,
		y:                 50,
	}

	return s, nil
}

// Append ...
func (s *DebugText) Append(str string) {
	s.text += str
}

// Update ...
func (s *DebugText) Update() error {
	l := s.uiFontColorLooper.Get()
	s.uiFontColor = color.RGBA{uint8(l), uint8(l), uint8(l), 255}

	return nil
}

// Draw ...
func (s *DebugText) Draw(r *ebiten.Image) {
	text.Draw(r, s.text, s.uiFont, s.x, s.y, s.uiFontColor)
	s.text = ""
}
