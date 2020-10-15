package layer

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten"
	"github.com/myanagisawa/ebitest/example/t5/ebitest"
	"github.com/myanagisawa/ebitest/example/t5/interfaces"
	"github.com/myanagisawa/ebitest/example/t5/models"
)

// Base ...
type Base struct {
	label    string
	bg       *models.EbiObject
	parent   interfaces.Scene
	isModal  bool
	controls []interfaces.UIControl
}

// Label ...
func (l *Base) Label() string {
	return fmt.Sprintf("%s.%s", l.parent.Label(), l.label)
}

// EbiObjects ...
func (l *Base) EbiObjects() []*models.EbiObject {
	return []*models.EbiObject{l.bg}
}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (l *Base) In(x, y int) bool {
	// レイヤ位置（左上座標）
	tx, ty := l.bg.GlobalTransition()

	// レイヤサイズ(オリジナル)
	w, h := l.bg.EbitenImage().Size()

	// スケール
	sx, sy := l.bg.GlobalScale()

	// 見かけ上の右下座標を取得
	maxX := int(float64(w)*sx + tx)
	maxY := int(float64(h)*sy + ty)
	if maxX > ebitest.Width {
		maxX = ebitest.Width
	}
	if maxY > ebitest.Height {
		maxY = ebitest.Height
	}

	// 見かけ上の左上座標を取得
	minX, minY := int(tx), int(ty)
	if minX < 0 {
		minX = 0
	}
	if minY < 0 {
		minY = 0
	}
	// log.Printf("レイヤ座標: {(%d, %d), (%d, %d)}", minX, minY, maxX, maxY)
	return (x >= minX && x <= maxX) && (y > minY && y <= maxY)
	// return l.bg.At(x-l.x, y-l.y).(color.RGBA).A > 0
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

// Update ...
func (l *Base) Update(screen *ebiten.Image) error {

	if l.parent.ActiveLayer() != nil && l.parent.ActiveLayer().Label() == l.label {
		// log.Printf("LayerBase.Update()")
		for _, c := range l.controls {

			_ = c.Update(screen)
		}
	}

	return nil
}

// Draw ...
func (l *Base) Draw(screen *ebiten.Image) {
	// log.Printf("LayerBase.Draw")
	op := &ebiten.DrawImageOptions{}

	w, h := l.bg.Size()
	// 描画位置指定
	op.GeoM.Reset()

	op.GeoM.Scale(l.bg.GlobalScale())

	// 対象画像の縦横半分だけマイナス位置に移動（原点に中心座標が来るように移動する）
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	// 中心を軸に回転
	op.GeoM.Rotate(l.bg.Theta())
	// ユニットの座標に移動
	op.GeoM.Translate(float64(w)/2, float64(h)/2)

	op.GeoM.Translate(l.bg.GlobalTransition())

	screen.DrawImage(l.bg.EbitenImage(), op)

	for _, c := range l.controls {
		c.Draw(screen)
	}

}

// NewLayerBase ...
func NewLayerBase(label string, img image.Image, parent interfaces.Scene, scale *ebitest.Scale, position *ebitest.Point, angle int) *Base {
	eimg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	l := &Base{
		label:  label,
		bg:     models.NewEbiObject(label, eimg, nil, scale, position, angle, false, false),
		parent: parent,
	}
	return l
}
