package layer

import (
	"image"
	"math"

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
	for _, route := range l.Routes {
		l.drawRoute(route, screen)
	}

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
	sitePos := l.getSitePosition(site)
	op.GeoM.Translate(sitePos.X(), sitePos.Y())

	// log.Printf("site: %s, (%0.2f, %0.2f)", site.Code, x, y)
	screen.DrawImage(site.Image, op)

	// site名を描画
	textSize := ebitest.NewSize(site.Text.Size())
	op.GeoM.Translate(-float64(textSize.W()-size.W())/2, float64(size.H()))
	screen.DrawImage(site.Text, op)
}

func (l *MapLayer) getSitePosition(site *parts.Site) *ebitest.Point {
	bgSize := l.bg.Size()
	x := l.bg.GlobalPosition().X() + float64(bgSize.W())*site.Location.X()
	y := l.bg.GlobalPosition().Y() + float64(bgSize.H())*site.Location.Y()

	return ebitest.NewPoint(x, y)
}

func (l *MapLayer) drawRoute(route *parts.Route, screen *ebiten.Image) {
	pos1 := l.getSitePosition(route.Site1)
	pos2 := l.getSitePosition(route.Site2)

	size1 := ebitest.NewSize(route.Site1.Image.Size())
	size2 := ebitest.NewSize(route.Site2.Image.Size())

	// ebitenutil.DrawLine(screen, pos1.X()+float64(size1.W())/2, pos1.Y()+float64(size1.H())/2, pos2.X()+float64(size2.W())/2, pos2.Y()+float64(size2.H())/2, color.Black)

	// 2点の座標
	x1, y1 := pos1.X()+float64(size1.W())/2, pos1.Y()+float64(size1.H())/2
	x2, y2 := pos2.X()+float64(size2.W())/2, pos2.Y()+float64(size2.H())/2

	// 距離を算出
	distance := math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))

	// 角度を算出
	rad := math.Atan2(y2-y1, x2-x1)
	// degree := rad * 180 / math.Pi

	op := &ebiten.DrawImageOptions{}

	routeSize := ebitest.NewSize(route.Image.Size())
	// 描画位置指定
	op.GeoM.Reset()
	op.GeoM.Scale(distance/float64(routeSize.W()), 1.0)

	op.GeoM.Rotate(rad)
	// log.Printf("degree: %0.2f, distance: %0.2f", degree, distance)
	// size :	= ebitest.NewSize(site.Image.Size())
	// // 対象画像の縦横半分だけマイナス位置に移動（原点に中心座標が来るように移動する）
	// op.GeoM.Translate(-float64(size.W())/2, -float64(size.H())/2)
	// // ユニットの座標に移動
	// op.GeoM.Translate(float64(size.W())/2, float64(size.H())/2)

	// // Site描画座標を計算
	// bgSize := l.bg.Size()
	// x := l.bg.GlobalPosition().X() + float64(bgSize.W())*site.Location.X()
	// y := l.bg.GlobalPosition().Y() + float64(bgSize.H())*site.Location.Y()
	op.GeoM.Translate(x1, y1)
	// // log.Printf("site: %s, (%0.2f, %0.2f)", site.Code, x, y)
	screen.DrawImage(route.Image, op)

}
