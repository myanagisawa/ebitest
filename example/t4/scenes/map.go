package scenes

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Map ...
type Map struct {
	LayerBase
	strokes map[*Stroke]struct{}
}

// NewBattleMap ...
func NewBattleMap(parent Scene) *Map {
	eimg := ebiten.NewImageFromImage(images["bgImage"])

	l := &Map{
		LayerBase: LayerBase{
			label:    "map",
			bg:       eimg,
			x:        0,
			y:        0,
			scale:    1.0,
			parent:   parent,
			isModal:  false,
			controls: []UIControl{},
			eventHandler: &EventHandler{
				events: map[string]map[*Event]struct{}{},
			},
		},
		strokes: map[*Stroke]struct{}{},
	}
	l.translateX = float64(l.x)
	l.translateY = float64(l.y)

	c := NewButton("メニューに戻る", images["btnBase"], fonts["btnFont"], color.Black, 100, 50)
	l.AddUIControl(c)
	l.AddEventListener(c, "click", func(target UIControl, source *EventSource) {
		log.Printf("メニューに戻る clicked: x=%d, y=%d", source.x, source.y)
		log.Printf("target: %#v", target)
		source.scene.Manager().TransitionToMainMenu()
	})

	c = NewButton("拡大", images["btnBaseHover"], fonts["btnFont"], color.White, 350, 50)
	l.AddUIControl(c)
	l.AddEventListener(c, "click", func(target UIControl, source *EventSource) {
		log.Printf("拡大 clicked: x=%d, y=%d", source.x, source.y)
		source.scene.ActiveLayer().ScaleTo(1.5)
	})

	c = NewButton("標準", images["btnBaseHover"], fonts["btnFont"], color.White, 600, 50)
	l.AddUIControl(c)
	l.AddEventListener(c, "click", func(target UIControl, source *EventSource) {
		log.Printf("標準 clicked: x=%d, y=%d", source.x, source.y)
		source.scene.ActiveLayer().ScaleTo(1.0)
	})

	c = NewButton("縮小", images["btnBaseHover"], fonts["btnFont"], color.White, 850, 50)
	l.AddUIControl(c)
	l.AddEventListener(c, "click", func(target UIControl, source *EventSource) {
		log.Printf("縮小 clicked: x=%d, y=%d", source.x, source.y)
		source.scene.ActiveLayer().ScaleTo(0.5)
	})

	return l
}

func (m *Map) updateStroke(stroke *Stroke) {
	stroke.Update()
	// if !stroke.IsReleased() {
	// 	return
	// }

	b := stroke.DraggingObject()
	if b == nil {
		// ドラッグ対象なしの場合はマップ自体のスクロール
		m.BgMoveBy(stroke.PositionDiff())
	}

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
	// log.Printf("MoveBy: s.x=%0.2f, s.y=%0.2f", m.translateX, m.translateY)
}

// Update ...
func (m *Map) Update() error {

	for stroke := range m.strokes {
		m.updateStroke(stroke)
		if stroke.IsReleased() {
			m.x, m.y = int(m.translateX), int(m.translateY)
			delete(m.strokes, stroke)
			log.Printf("drag end")
		}
	}

	if m.parent.ActiveLayer() == m {
		x, y := m.LocalPosition(ebiten.CursorPosition())
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			stroke := NewStroke(&MouseStrokeSource{})
			// レイヤ内のドラッグ対象のオブジェクトを取得する仕組みが必要
			stroke.SetDraggingObject(m.UIControlAt(ebiten.CursorPosition()))
			m.strokes[stroke] = struct{}{}
			log.Printf("drag start")
		}

		for _, c := range m.controls {
			switch val := c.(type) {
			case *UIButtonImpl:
				val.hover = val.In(x, y)
			default:
				log.Printf("Map.Update: controls switch: default: %#v", val)
			}
			// _ = c.Update(screen)
		}
	}

	m.LayerBase.Update()
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

	for _, c := range m.controls {
		c.Draw(m.bg)
	}
}
