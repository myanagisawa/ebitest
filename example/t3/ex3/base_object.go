package ex3

import (
	"image/color"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

type (
	// Circle ...
	Circle struct {
		r     int
		image ebiten.Image
	}

	// Scene ...
	Scene interface {
		Update() error
		Draw(r *ebiten.Image)
		GetSize() (int, int)
	}

	// SceneImpl ...
	SceneImpl struct {
	}

	// Point ...
	Point struct {
		X int
		Y int
	}

	// Size ...
	Size struct {
		Width  int
		Height int
	}

	// Looper ...
	Looper struct {
		num int
		vec int
		min int
		max int
	}

	// LabelFace ...
	LabelFace struct {
		uiFont        font.Face
		uiFontColor   color.Color
		uiFontMHeight int
	}
)

// NewLabelFace ...
func NewLabelFace(size int, c color.Color) *LabelFace {
	// ebitenフォントのテスト
	tt, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	uiFont := truetype.NewFace(tt, &truetype.Options{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	b, _, _ := uiFont.GlyphBounds('M')
	uiFontMHeight := (b.Max.Y - b.Min.Y).Ceil()

	s := &LabelFace{
		uiFont:        uiFont,
		uiFontColor:   c,
		uiFontMHeight: uiFontMHeight,
	}
	return s
}

// NewSceneImpl ...
func NewSceneImpl() *SceneImpl {
	s := &SceneImpl{}
	return s
}

// Update ...
func (s *SceneImpl) Update() error {
	return nil
}

// Draw ...
func (s *SceneImpl) Draw(r *ebiten.Image) {
}

// GetSize ...
func (s *SceneImpl) GetSize() (int, int) {
	return 0, 0
}

// NewLooper ...
func NewLooper(num, step, min, max int) *Looper {
	l := &Looper{num, step, min, max}
	return l
}

// Get ...
func (l *Looper) Get() int {
	ret := l.num
	if l.num <= l.min || l.num >= l.max {
		l.vec = l.vec * -1
	}
	l.num += l.vec
	return ret
}

// R ...
func (c *Circle) R() int {
	return c.r
}
