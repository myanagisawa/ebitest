package layer

import (
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/myanagisawa/ebitest/example/t5/ebitest"
	"github.com/myanagisawa/ebitest/example/t5/enum"
	"github.com/myanagisawa/ebitest/example/t5/interfaces"
	"github.com/myanagisawa/ebitest/example/t5/models"
	"github.com/myanagisawa/ebitest/example/t5/models/event"
	"github.com/myanagisawa/ebitest/example/t5/models/input"
)

// Base ...
type Base struct {
	label        string
	bg           *models.EbiObject
	parent       interfaces.Scene
	isModal      bool
	controls     []interfaces.UIControl
	eventHandler *event.Handler
	stroke       *input.Stroke
}

// Label ...
func (l *Base) Label() string {
	return l.label
}

// LabelFull ...
func (l *Base) LabelFull() string {
	return fmt.Sprintf("%s.%s", l.parent.Label(), l.label)
}

// EbiObjects ...
func (l *Base) EbiObjects() []*models.EbiObject {
	return []*models.EbiObject{l.bg}
}

// EventHandler ...
func (l *Base) EventHandler() interfaces.EventHandler {
	return l.eventHandler
}

// Scroll ...
func (l *Base) Scroll(et enum.EdgeTypeEnum) {
	// log.Printf("%s: EdgeType: %d", l.Label(), et)
	// 1フレームあたりの増分値
	dp := 20.0
	switch et {
	case enum.EdgeTypeTopLeft:
		l.bg.Position().SetDelta(dp, dp)
	case enum.EdgeTypeTop:
		l.bg.Position().SetDelta(0, dp)
	case enum.EdgeTypeTopRight:
		l.bg.Position().SetDelta(-dp, dp)
	case enum.EdgeTypeRight:
		l.bg.Position().SetDelta(-dp, 0)
	case enum.EdgeTypeBottomRight:
		l.bg.Position().SetDelta(-dp, -dp)
	case enum.EdgeTypeBottom:
		l.bg.Position().SetDelta(0, -dp)
	case enum.EdgeTypeBottomLeft:
		l.bg.Position().SetDelta(dp, -dp)
	case enum.EdgeTypeLeft:
		l.bg.Position().SetDelta(dp, 0)
	}

	bgPos := l.bg.GlobalPosition()
	bgSize := l.bg.Size()
	// log.Printf("global position: %0.1f, %0.1f", gx, gy)
	if int(bgPos.X())+bgSize.W() < ebitest.Width {
		l.bg.Position().SetDelta(dp, 0)
	} else if bgPos.X() > 0 {
		l.bg.Position().SetDelta(-dp, 0)
	}
	if int(bgPos.Y())+bgSize.H() < ebitest.Height {
		l.bg.Position().SetDelta(0, dp)
	} else if bgPos.Y() > 0 {
		l.bg.Position().SetDelta(0, -dp)
	}
}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (l *Base) In(x, y int) bool {
	// // レイヤ位置（左上座標）
	// tx, ty := l.bg.GlobalPosition()

	// // レイヤサイズ(オリジナル)
	// w, h := l.bg.EbitenImage().Size()

	// // スケール
	// sx, sy := l.bg.GlobalScale()

	// // 見かけ上の右下座標を取得
	// maxX := int(float64(w)*sx + tx)
	// maxY := int(float64(h)*sy + ty)
	// if maxX > ebitest.Width {
	// 	maxX = ebitest.Width
	// }
	// if maxY > ebitest.Height {
	// 	maxY = ebitest.Height
	// }

	// // 見かけ上の左上座標を取得
	// minX, minY := int(tx), int(ty)
	// if minX < 0 {
	// 	minX = 0
	// }
	// if minY < 0 {
	// 	minY = 0
	// }
	// // log.Printf("カーソル位置: (%d, %d)  レイヤ座標: {(%d, %d), (%d, %d)}", x, y, minX, minY, maxX, maxY)
	// return (x >= minX && x <= maxX) && (y > minY && y <= maxY)
	// // return l.bg.At(x-l.x, y-l.y).(color.RGBA).A > 0
	return l.bg.In(x, y)
}

// IsModal ...
func (l *Base) IsModal() bool {
	return l.isModal
}

// AddUIControl レイヤに部品を追加します
func (l *Base) AddUIControl(c interfaces.UIControl) {
	c.SetLayer(l)
	l.controls = append(l.controls, c)
}

// UIControlAt (x, y)座標に存在する部品を返します
func (l *Base) UIControlAt(x, y int) interfaces.UIControl {
	for i := len(l.controls) - 1; i >= 0; i-- {
		c := l.controls[i]
		if c.In(x, y) {
			return c
		}
	}
	return nil
}

func (l *Base) updateStroke(stroke *input.Stroke) {
	stroke.Update()
	// if !stroke.IsReleased() {
	// 	return
	// }
	l.bg.SetMoving(stroke.PositionDiff())
}

// Update ...
func (l *Base) Update() error {

	if l.stroke != nil {
		l.updateStroke(l.stroke)
		if l.stroke.IsReleased() {
			l.bg.UpdatePositionByDelta()
			l.stroke = nil
			log.Printf("drag end")
		}
	}

	if l.parent.ActiveLayer() != nil && l.parent.ActiveLayer() == l {
		// log.Printf("active layer: %s", l.parent.ActiveLayer().Label())

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			// log.Printf("left button push: x=%d, y=%d", x, y)
			if l.In(x, y) {
				stroke := input.NewStroke(&input.MouseStrokeSource{})
				// レイヤ内のドラッグ対象のオブジェクトを取得する仕組みが必要
				// o := l.UIControlAt(x, y)
				// if o != nil || l.bg.IsDraggable() {
				// 	l.stroke = stroke
				// 	log.Printf("drag start")
				// }
				if l.bg.IsDraggable() {
					l.stroke = stroke
					log.Printf("%s drag start", l.label)
				}
			}
		}

		// log.Printf("LayerBase.Update()")
		for _, c := range l.controls {
			_ = c.Update()
		}
	}

	return nil
}

// Draw ...
func (l *Base) Draw(screen *ebiten.Image) {
	// log.Printf("LayerBase.Draw")
	op := &ebiten.DrawImageOptions{}

	// 描画位置指定
	op.GeoM.Reset()
	op.GeoM.Scale(l.bg.GlobalScale().Get())

	bgSize := l.bg.Size()
	// 対象画像の縦横半分だけマイナス位置に移動（原点に中心座標が来るように移動する）
	op.GeoM.Translate(-float64(bgSize.W())/2, -float64(bgSize.H())/2)
	// 中心を軸に回転
	op.GeoM.Rotate(l.bg.Theta())
	// ユニットの座標に移動
	op.GeoM.Translate(float64(bgSize.W())/2, float64(bgSize.H())/2)

	op.GeoM.Translate(l.bg.GlobalPosition().Get())

	screen.DrawImage(l.bg.EbitenImage(), op)

	for _, c := range l.controls {
		c.Draw(screen)
	}

}

// NewLayerBase ...
func NewLayerBase(label string, img image.Image, parent interfaces.Scene, scale *ebitest.Scale, position *ebitest.Point, angle int, draggable bool) *Base {
	eimg := ebiten.NewImageFromImage(img)

	l := &Base{
		label:        label,
		bg:           models.NewEbiObject(label, eimg, nil, scale, position, angle, false, false, draggable),
		parent:       parent,
		eventHandler: event.NewEventHandler(),
	}
	return l
}
