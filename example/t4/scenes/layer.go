package scenes

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
)

// Layer ...
type Layer interface {
	In(x, y int) bool
	IsModal() bool
	Update(screen *ebiten.Image) error
	Draw(screen *ebiten.Image)
}

// LayerBase ...
type LayerBase struct {
	bg         *ebiten.Image
	x          int
	y          int
	scale      float64
	translateX float64
	translateY float64
	parent     *BattleScene
	isModal    bool
}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (l *LayerBase) In(x, y int) bool {
	return l.bg.At(x-l.x, y-l.y).(color.RGBA).A > 0
}

// IsModal ...
func (l *LayerBase) IsModal() bool {
	return l.isModal
}

// Update ...
func (l *LayerBase) Update(screen *ebiten.Image) error {

	return nil
}

// Draw ...
func (l *LayerBase) Draw(screen *ebiten.Image) {
	log.Printf("LayerBase.Draw")
	screen.Fill(color.RGBA{200, 200, 200, 64})
}

// TestWindow ...
type TestWindow struct {
	LayerBase
}

// NewTestWindow ...
func NewTestWindow() *TestWindow {
	img := createRectImage(200, 400, color.RGBA{0, 0, 0, 64})
	eimg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	l := &TestWindow{
		LayerBase: LayerBase{
			bg:      eimg,
			x:       0,
			y:       0,
			scale:   1.0,
			isModal: false,
		},
	}
	return l
}

// Draw ...
func (l *TestWindow) Draw(screen *ebiten.Image) {
	log.Printf("TestWindow.Draw")

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(l.scale, l.translateY)

	screen.DrawImage(l.bg, op)
}
