package layer

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/example/t5/ebitest"
	"github.com/myanagisawa/ebitest/example/t5/interfaces"
	"github.com/myanagisawa/ebitest/example/t5/parts"
)

// MapLayer ...
type MapLayer struct {
	Base
	Sites  []*parts.Site
	Routes []*parts.Route
}

// NewMapLayer ...
func NewMapLayer(label string, img image.Image, parent interfaces.Scene, scale *ebitest.Scale, position *ebitest.Point, angle int, draggable bool) *MapLayer {
	baseLayer := NewLayerBase(label, img, parent, scale, position, angle, draggable)
	return &MapLayer{
		Base: *baseLayer,
	}
}

// Draw ...
func (l *MapLayer) Draw(screen *ebiten.Image) {
	// 共通Draw処理
	l.Base.Draw(screen)

	// log.Printf("---------------- site draw --------------")
	for _, site := range l.Sites {
		l.drawSite(site, screen)
	}
}

func (l *MapLayer) drawSite(site *parts.Site, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	// 描画位置指定
	op.GeoM.Reset()
	op.GeoM.Scale(1.0, 1.0)

	size := ebitest.NewSize(site.Image.Size())
	// 対象画像の縦横半分だけマイナス位置に移動（原点に中心座標が来るように移動する）
	op.GeoM.Translate(-float64(size.W())/2, -float64(size.H())/2)
	// ユニットの座標に移動
	op.GeoM.Translate(float64(size.W())/2, float64(size.H())/2)

	// Site描画座標を計算
	bgSize := l.bg.Size()
	x := l.bg.GlobalPosition().X() + float64(bgSize.W())*site.Location.X()
	y := l.bg.GlobalPosition().Y() + float64(bgSize.H())*site.Location.Y()
	op.GeoM.Translate(x, y)

	// log.Printf("site: %s, (%0.2f, %0.2f)", site.Code, x, y)
	screen.DrawImage(site.Image, op)
}
