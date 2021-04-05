package m

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/app/g"
	"github.com/myanagisawa/ebitest/example/t7/lib/models/char"
	"github.com/myanagisawa/ebitest/example/t7/lib/utils"
)

// Route マップ上の経路情報
type Route struct {
	Code  string
	Type  enum.RouteTypeEnum
	Name  string
	Site1 *Site
	Site2 *Site
	Image image.Image
	Text  image.Image
}

// NewRoute ...
func NewRoute(code string, t enum.RouteTypeEnum, name string, site1, site2 *Site) *Route {
	if code == "" || site1 == nil || site2 == nil {
		panic(fmt.Sprintf("Route初期化エラー: code=%#v, site1=%#v, site2=%#v", code, site1, site2))
	}

	fset := char.Res.Get(12, enum.FontStyleGenShinGothicRegular)
	ti := fset.GetStringImage(name)
	ti = utils.TextColorTo(ti.(draw.Image), &color.RGBA{255, 255, 255, 255})

	return &Route{
		Code:  code,
		Type:  t,
		Name:  name,
		Site1: site1,
		Site2: site2,
		Image: g.Images[fmt.Sprintf("route_%d", t)],
		Text:  ti,
	}
}

// Routes ...
type Routes []Route

// GetByCode ...
func (o *Routes) GetByCode(code string) *Route {
	for _, r := range *o {
		if r.Code == code {
			return &r
		}
	}
	return nil
}

// CreateRoutes ...
func CreateRoutes(rawRoutes []RawRoute, sites *Sites) *Routes {
	routes := Routes{}
	for _, route := range rawRoutes {
		s1 := sites.GetByCode(route.Site1)
		s2 := sites.GetByCode(route.Site2)
		if s1 == nil || s2 == nil {
			log.Printf("invalid type: ")
			continue
		}

		routes = append(routes, *NewRoute(route.Code, route.Type, route.Name, s1, s2))
	}
	return &routes
}
