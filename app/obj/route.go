package obj

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/app/g"
	"github.com/myanagisawa/ebitest/enum"
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

	img := ebiten.NewImageFromImage(g.Images[fmt.Sprintf("route_%d", t)])

	return &Route{
		Code:  code,
		Type:  t,
		Name:  name,
		Site1: site1,
		Site2: site2,
		Image: img,
	}
}

// Routes ...
type Routes []Route

// CreateRoutes ...
func CreateRoutes(routeMaster []g.MRoute, sites Sites) Routes {
	routes := Routes{}
	for _, route := range routeMaster {
		routes = append(routes, *NewRoute(route.Code, route.Type, route.Name, sites.GetSiteByCode(route.Site1), sites.GetSiteByCode(route.Site2)))
	}
	return routes
}
