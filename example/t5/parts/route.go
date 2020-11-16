package parts

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/example/t5/ebitest"
	"github.com/myanagisawa/ebitest/example/t5/enum"
)

// Route マップ上の経路情報
type Route struct {
	Code  string
	Type  enum.RouteTypeEnum
	Name  string
	Site1 *Site
	Site2 *Site
	Image *ebiten.Image
}

// NewRoute ...
func NewRoute(code string, t enum.RouteTypeEnum, name string, site1, site2 *Site) *Route {
	if code == "" || site1 == nil || site2 == nil {
		panic(fmt.Sprintf("Route初期化エラー: code=%#v, site1=%#v, site2=%#v", code, site1, site2))
	}

	img := ebiten.NewImageFromImage(ebitest.Images[fmt.Sprintf("route_%d", t)])

	return &Route{
		Code:  code,
		Type:  t,
		Name:  name,
		Site1: site1,
		Site2: site2,
		Image: img,
	}
}
