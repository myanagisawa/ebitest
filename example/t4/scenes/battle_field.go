package scenes

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type (

	// BattleScene ...
	BattleScene struct {
		manager *GameManager
		field   *Map
		layers  []Layer
		strokes map[*Stroke]struct{}
	}
)

func (s *BattleScene) updateStroke(stroke *Stroke) {
	stroke.Update()
	// if !stroke.IsReleased() {
	// 	return
	// }

	b := stroke.DraggingObject()
	if b == nil {
		return
	}

	s.field.BgMoveBy(stroke.PositionDiff())

	// stroke.SetDraggingObject(nil)
}

// NewBattleScene ...
func NewBattleScene(m *GameManager) Scene {
	s := &BattleScene{
		manager: m,
		strokes: map[*Stroke]struct{}{},
	}

	l1 := NewTestWindow()
	l1.parent = s
	s.layers = append(s.layers, l1)

	s.field = NewBattleMap(s)

	return s
}

// GetActiveLayer ...
func (s *BattleScene) GetActiveLayer() *Layer {
	for _, layer := range s.layers {
		if layer.IsModal() {
			return &layer
		}
		if layer.In(ebiten.CursorPosition()) {
			return &layer
		}
	}
	return nil
}

// Update ...
func (s *BattleScene) Update(screen *ebiten.Image) error {

	activeLayer := s.GetActiveLayer()
	log.Printf("activeLayer: %#v", activeLayer)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		stroke := NewStroke(&MouseStrokeSource{})
		// レイヤ内のドラッグ対象のオブジェクトを取得する仕組みが必要
		stroke.SetDraggingObject(s.field.bg)
		s.strokes[stroke] = struct{}{}
		log.Printf("drag start")
	}

	for stroke := range s.strokes {
		s.updateStroke(stroke)
		if stroke.IsReleased() {
			s.field.x, s.field.y = int(s.field.translateX), int(s.field.translateY)
			delete(s.strokes, stroke)
			log.Printf("drag end")
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		// click イベントを発火
		x, y := s.field.LocalPosition(ebiten.CursorPosition())
		eventHandler.Firing(s, "click", x, y)
	}

	s.field.Update(screen)

	for _, layer := range s.layers {
		layer.Update(screen)
	}

	return nil
}

// Draw ...
func (s *BattleScene) Draw(screen *ebiten.Image) {

	s.field.Draw(screen)

	for _, layer := range s.layers {
		layer.Draw(screen)
	}
}

// Layout ...
func (s *BattleScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}
