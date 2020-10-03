package scenes

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
)

// Map ...
type Map struct {
	bg          *ebiten.Image
	x           int
	y           int
	scale       float64
	translateX  float64
	translateY  float64
	parent      *BattleScene
	controllers map[UIController]struct{}
}

// NewBattleMap ...
func NewBattleMap(parent *BattleScene) *Map {
	eimg, _ := ebiten.NewImageFromImage(images["bgImage"], ebiten.FilterDefault)

	f := &Map{
		bg:          eimg,
		x:           0,
		y:           0,
		scale:       1.0,
		parent:      parent,
		controllers: map[UIController]struct{}{},
	}
	f.translateX = float64(f.x)
	f.translateY = float64(f.y)

	c := NewButton("Return to menu", images["btnBase"], fonts["btnFont"], color.Black, 100, 200)
	c.AddEventListener("click", func(target UIController, source *EventSource) {
		log.Printf("btnRtn clicked")
		s := source.scene.(*BattleScene)
		s.manager.TransitionToMainMenu()
	})
	f.controllers[c] = struct{}{}

	c = NewButton("Scale to 1.5", images["btnBaseHover"], fonts["btnFont"], color.White, 100, 300)
	c.AddEventListener("click", func(target UIController, source *EventSource) {
		log.Printf("btnRtn clicked")
		s := source.scene.(*BattleScene)
		s.field.ScaleTo(1.5)
	})
	f.controllers[c] = struct{}{}

	c = NewButton("Scale to 1.0", images["btnBaseHover"], fonts["btnFont"], color.White, 100, 400)
	c.AddEventListener("click", func(target UIController, source *EventSource) {
		log.Printf("btnRtn clicked")
		s := source.scene.(*BattleScene)
		s.field.ScaleTo(1.0)
	})
	f.controllers[c] = struct{}{}

	c = NewButton("Scale to 0.5", images["btnBaseHover"], fonts["btnFont"], color.White, 100, 500)
	c.AddEventListener("click", func(target UIController, source *EventSource) {
		log.Printf("btnRtn clicked")
		s := source.scene.(*BattleScene)
		s.field.ScaleTo(0.5)
	})
	f.controllers[c] = struct{}{}

	return f
}

// BgMoveBy moves the background by (x, y).
func (m *Map) BgMoveBy(x, y int) {
	w, h := m.bg.Size()
	w = int(float64(w) * m.scale)
	h = int(float64(h) * m.scale)

	m.translateX = float64(m.x + x)
	m.translateY = float64(m.y + y)
	if m.translateX > 0 {
		m.translateX = 0
	}
	if m.translateX < float64(width-w) {
		m.translateX = float64(width - w)
	}
	// log.Printf("w(%d) - transX(%d) < width(%d): %v", w, int(math.Abs(s.translateX*s.scale)), width, (w-int(math.Abs(s.translateX))) < width)
	if (w - int(math.Abs(m.translateX*m.scale))) < width {
		m.translateX = float64(width-w) / m.scale
	}
	if m.translateY > 0 {
		m.translateY = 0
	}
	if m.translateY < float64(height-h) {
		m.translateY = float64(height - h)
	}
	if (h - int(math.Abs(m.translateY*m.scale))) < height {
		m.translateY = float64(height-h) / m.scale
	}
	// log.Printf("MoveBy: s.x=%0.2f, s.y=%0.2f", s.translateX, s.translateY)
}

// ScaleTo ...
func (m *Map) ScaleTo(scale float64) {
	m.scale = scale

	w, h := m.bg.Size()
	w = int(float64(w) * m.scale)
	h = int(float64(h) * m.scale)

	if m.x < width-w {
		m.x = width - w
	}
	if m.y < height-h {
		m.y = height - h
	}
	m.translateX = float64(m.x)
	m.translateY = float64(m.y)
	// log.Printf("MoveBy: s.x=%0.2f, s.y=%0.2f", s.translateX, s.translateY)
}

// LocalPosition スクリーン上の座標をシーンオブジェクト上の座標に変換します
func (m *Map) LocalPosition(x, y int) (int, int) {
	sx := float64(m.x) * m.scale * -1
	sy := float64(m.y) * m.scale * -1
	localX := int((float64(x) + sx) / m.scale)
	localY := int((float64(y) + sy) / m.scale)

	// log.Printf("scale: %0.2f [x: %d, s.x: %d = %d] [y: %d, s.y: %d = %d]", s.scale, x, s.x, localX, y, s.y, localY)
	return localX, localY
}

// Update ...
func (m *Map) Update(screen *ebiten.Image) error {

	for c := range m.controllers {
		switch val := c.(type) {
		case *UIButtonImpl:
			val.hover = val.In(m.LocalPosition(ebiten.CursorPosition()))
		default:
			log.Printf("Map.Update: controllers switch: default: %#v", val)
		}
		_ = c.Update(screen)
	}

	// log.Printf("bg.x=%d, bg.y=%d", s.x, s.y)
	return nil
}

// Draw ...
func (m *Map) Draw(screen *ebiten.Image) {

	// Draw background
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(0.5, 0.5, 0.5, 1.0)
	op.GeoM.Translate(m.translateX, m.translateY)
	op.GeoM.Scale(m.scale, m.scale)

	screen.DrawImage(m.bg, op)
	// screen.DrawImage(s.bg.SubImage(image.Rect(300, 200, 1500, 1100)).(*ebiten.Image), op)

	for c := range m.controllers {
		c.Draw(m.bg)
	}
}
