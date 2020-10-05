package scenes

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// Map ...
type Map struct {
	LayerBase
	strokes map[*Stroke]struct{}
}

// NewBattleMap ...
func NewBattleMap(parent Scene) *Map {
	eimg, _ := ebiten.NewImageFromImage(images["bgImage"], ebiten.FilterDefault)

	l := &Map{
		LayerBase: LayerBase{
			label:    "map",
			bg:       eimg,
			x:        0,
			y:        0,
			scale:    1.0,
			parent:   parent,
			isModal:  false,
			controls: map[UIController]struct{}{},
		},
		strokes: map[*Stroke]struct{}{},
	}
	l.translateX = float64(l.x)
	l.translateY = float64(l.y)

	c := NewButton("Return to menu", images["btnBase"], fonts["btnFont"], color.Black, 100, 50)
	c.AddEventListener(parent, "click", func(target UIController, source *EventSource) {
		log.Printf("btnRtn clicked")
		source.scene.Manager().TransitionToMainMenu()
	})
	l.controls[c] = struct{}{}

	c = NewButton("Scale to 1.5", images["btnBaseHover"], fonts["btnFont"], color.White, 350, 50)
	c.AddEventListener(parent, "click", func(target UIController, source *EventSource) {
		log.Printf("btnRtn clicked")
		source.scene.ActiveLayer().ScaleTo(1.5)
	})
	l.controls[c] = struct{}{}

	c = NewButton("Scale to 1.0", images["btnBaseHover"], fonts["btnFont"], color.White, 600, 50)
	c.AddEventListener(parent, "click", func(target UIController, source *EventSource) {
		log.Printf("btnRtn clicked")
		source.scene.ActiveLayer().ScaleTo(1.0)
	})
	l.controls[c] = struct{}{}

	c = NewButton("Scale to 0.5", images["btnBaseHover"], fonts["btnFont"], color.White, 850, 50)
	c.AddEventListener(parent, "click", func(target UIController, source *EventSource) {
		log.Printf("btnRtn clicked")
		source.scene.ActiveLayer().ScaleTo(0.5)
	})
	l.controls[c] = struct{}{}

	return l
}

func (m *Map) updateStroke(stroke *Stroke) {
	stroke.Update()
	// if !stroke.IsReleased() {
	// 	return
	// }

	b := stroke.DraggingObject()
	if b == nil {
		return
	}

	m.BgMoveBy(stroke.PositionDiff())

	// stroke.SetDraggingObject(nil)
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

// Update ...
func (m *Map) Update(screen *ebiten.Image) error {
	if m.parent.ActiveLayer() == m {
		x, y := m.LocalPosition(ebiten.CursorPosition())
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			stroke := NewStroke(&MouseStrokeSource{})
			// レイヤ内のドラッグ対象のオブジェクトを取得する仕組みが必要
			stroke.SetDraggingObject(m.ControlAt(ebiten.CursorPosition()))
			m.strokes[stroke] = struct{}{}
			log.Printf("drag start")
		}

		for stroke := range m.strokes {
			m.updateStroke(stroke)
			if stroke.IsReleased() {
				m.x, m.y = int(m.translateX), int(m.translateY)
				delete(m.strokes, stroke)
				log.Printf("drag end")
			}
		}

		for c := range m.controls {
			switch val := c.(type) {
			case *UIButtonImpl:
				val.hover = val.In(x, y)
			default:
				log.Printf("Map.Update: controls switch: default: %#v", val)
			}
			_ = c.Update(screen)
		}
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

	for c := range m.controls {
		c.Draw(m.bg)
	}
}
