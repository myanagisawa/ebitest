package scenes

import (
	"image/color"
	"image/draw"

	"github.com/hajimehoshi/ebiten"
	"github.com/myanagisawa/ebitest/utils"
	"golang.org/x/image/font"
)

// UIControl ...
type UIControl interface {
	Update(screen *ebiten.Image) error
	Draw(screen *ebiten.Image)
	In(x, y int) bool
	SetLayer(l Layer)
}

// UIButton ...
type UIButton interface {
	UIControl
}

// UIControllerImpl ...
type UIControllerImpl struct {
	layer Layer
	image *ebiten.Image
	x     int
	y     int
}

// SetLayer ...
func (c *UIControllerImpl) SetLayer(l Layer) {
	c.layer = l
}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (c *UIControllerImpl) In(x, y int) bool {
	return c.image.At(x-c.x, y-c.y).(color.RGBA).A > 0
}

// Update ...
func (c *UIControllerImpl) Update(screen *ebiten.Image) error {
	return nil
}

// Draw draws the sprite.
func (c *UIControllerImpl) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(c.image, op)
}

// UIButtonImpl ...
type UIButtonImpl struct {
	UIControllerImpl
	hover bool
}

// NewButton ...
func NewButton(label string, baseImg draw.Image, fontFace font.Face, labelColor color.Color, x, y int) UIButton {
	img := utils.SetTextToCenter(label, baseImg, fontFace, labelColor)
	eimg, _ := ebiten.NewImageFromImage(*img, ebiten.FilterDefault)
	con := &UIControllerImpl{image: eimg, x: x, y: y}
	return &UIButtonImpl{UIControllerImpl: *con}
}

// Draw draws the sprite.
func (c *UIButtonImpl) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.x), float64(c.y))
	r, g, b, a := 1.0, 1.0, 1.0, 1.0
	if c.hover {
		r, g, b, a = 0.5, 0.5, 0.5, 1.0
	}
	op.ColorM.Scale(r, g, b, a)
	screen.DrawImage(c.image, op)
}
