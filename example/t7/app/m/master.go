package m

import (
	"fmt"

	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/app/g"
)

// MasterData ...
type MasterData struct {
	RawSites  []RawSite
	RawRoutes []RawRoute
}

// LoadMaster ...
func LoadMaster() *MasterData {
	m := &MasterData{}
	m.RawSites = LoadSites()
	m.RawRoutes = LoadRoutes()
	return m
}

// RawSite マップ上の拠点情報
type RawSite struct {
	Code     string
	Type     enum.SiteTypeEnum
	Name     string
	Location *g.Point
}

// NewRawSite ...
func NewRawSite(code string, t enum.SiteTypeEnum, name string, loc *g.Point) *RawSite {
	if code == "" || loc.X() > 1 || loc.Y() > 1 {
		panic(fmt.Sprintf("Site初期化エラー: code=%#v, loc=%#v", code, loc))
	}

	return &RawSite{
		Code:     code,
		Type:     t,
		Name:     name,
		Location: loc,
	}
}

// RawRoute マップ上の経路情報
type RawRoute struct {
	Code  string
	Type  enum.RouteTypeEnum
	Name  string
	Site1 string
	Site2 string
}

// NewRawRoute ...
func NewRawRoute(code string, t enum.RouteTypeEnum, name, site1, site2 string) *RawRoute {
	if code == "" || site1 == "" || site2 == "" {
		panic(fmt.Sprintf("Route初期化エラー: code=%#v, site1=%#v, site2=%#v", code, site1, site2))
	}

	return &RawRoute{
		Code:  code,
		Type:  t,
		Name:  name,
		Site1: site1,
		Site2: site2,
	}
}

// LoadSites ...
func LoadSites() []RawSite {
	return []RawSite{
		*NewRawSite("site-1", enum.SiteTypeBase, "base-1", g.NewPoint(0.318, 0.707)),
		*NewRawSite("site-2", enum.SiteTypeBase, "base-2", g.NewPoint(0.161, 0.761)),
		*NewRawSite("site-3", enum.SiteTypeBase, "base-3", g.NewPoint(0.237, 0.826)),
		*NewRawSite("site-4", enum.SiteTypeBase, "base-4", g.NewPoint(0.345, 0.847)),
		*NewRawSite("site-5", enum.SiteTypeBase, "base-5", g.NewPoint(0.255, 0.883)),
		*NewRawSite("site-6", enum.SiteTypeBase, "base-6", g.NewPoint(0.573, 0.780)),
		*NewRawSite("site-7", enum.SiteTypeCity, "city-1", g.NewPoint(0.306, 0.623)),
		*NewRawSite("site-8", enum.SiteTypeCity, "city-2", g.NewPoint(0.302, 0.663)),
		*NewRawSite("site-9", enum.SiteTypeCity, "city-3", g.NewPoint(0.295, 0.726)),
		*NewRawSite("site-10", enum.SiteTypeCity, "city-4", g.NewPoint(0.276, 0.751)),
		*NewRawSite("site-11", enum.SiteTypeCity, "city-5", g.NewPoint(0.181, 0.648)),
		*NewRawSite("site-12", enum.SiteTypeCity, "city-6", g.NewPoint(0.181, 0.598)),
		*NewRawSite("site-13", enum.SiteTypeCity, "city-7", g.NewPoint(0.120, 0.847)),
		*NewRawSite("site-14", enum.SiteTypeCity, "city-8", g.NewPoint(0.193, 0.863)),
		*NewRawSite("site-15", enum.SiteTypeCity, "city-9", g.NewPoint(0.263, 0.953)),
		*NewRawSite("site-16", enum.SiteTypeCity, "city-10", g.NewPoint(0.279, 0.842)),
		*NewRawSite("site-17", enum.SiteTypeCity, "city-11", g.NewPoint(0.314, 0.893)),
		*NewRawSite("site-18", enum.SiteTypeCity, "city-12", g.NewPoint(0.374, 0.835)),
		*NewRawSite("site-19", enum.SiteTypeCity, "city-13", g.NewPoint(0.378, 0.932)),
		*NewRawSite("site-20", enum.SiteTypeCity, "city-14", g.NewPoint(0.448, 0.863)),
		*NewRawSite("site-21", enum.SiteTypeCity, "city-15", g.NewPoint(0.479, 0.824)),
		*NewRawSite("site-22", enum.SiteTypeCity, "city-16", g.NewPoint(0.563, 0.852)),
		*NewRawSite("site-23", enum.SiteTypeCity, "city-17", g.NewPoint(0.565, 0.708)),
		*NewRawSite("site-24", enum.SiteTypeCity, "city-18", g.NewPoint(0.600, 0.613)),
		*NewRawSite("site-25", enum.SiteTypeCity, "city-19", g.NewPoint(0.652, 0.856)),
		*NewRawSite("site-26", enum.SiteTypeCity, "city-20", g.NewPoint(0.698, 0.805)),
		*NewRawSite("site-27", enum.SiteTypeFort, "fort-1", g.NewPoint(0.181, 0.718)),
		*NewRawSite("site-28", enum.SiteTypeFort, "fort-2", g.NewPoint(0.189, 0.785)),
		*NewRawSite("site-29", enum.SiteTypeFort, "fort-3", g.NewPoint(0.252, 0.801)),
		*NewRawSite("site-30", enum.SiteTypeFort, "fort-4", g.NewPoint(0.413, 0.834)),
		*NewRawSite("site-31", enum.SiteTypePoint, "point-1", g.NewPoint(0.245, 0.714)),
		*NewRawSite("site-32", enum.SiteTypePoint, "point-2", g.NewPoint(0.209, 0.822)),
		*NewRawSite("site-33", enum.SiteTypePoint, "point-3", g.NewPoint(0.230, 0.868)),
		*NewRawSite("site-34", enum.SiteTypePoint, "point-4", g.NewPoint(0.302, 0.807)),
		*NewRawSite("site-35", enum.SiteTypePoint, "point-5", g.NewPoint(0.484, 0.918)),
		*NewRawSite("site-36", enum.SiteTypePoint, "point-6", g.NewPoint(0.566, 0.751)),
		*NewRawSite("site-37", enum.SiteTypePoint, "point-7", g.NewPoint(0.580, 0.656)),
		*NewRawSite("site-38", enum.SiteTypeRemains, "remains-1", g.NewPoint(0.127, 0.730)),
		*NewRawSite("site-39", enum.SiteTypeRemains, "remains-2", g.NewPoint(0.223, 0.786)),
		*NewRawSite("site-40", enum.SiteTypeRemains, "remains-3", g.NewPoint(0.580, 0.747)),
		*NewRawSite("site-41", enum.SiteTypeRemains, "remains-4", g.NewPoint(0.636, 0.793)),
	}
}

// LoadRoutes ...
func LoadRoutes() []RawRoute {
	return []RawRoute{
		*NewRawRoute("route-1", enum.RouteTypeSea, "route-1", "site-11", "site-12"),
		*NewRawRoute("route-2", enum.RouteTypePlain, "route-2", "site-11", "site-27"),
		*NewRawRoute("route-3", enum.RouteTypeMountain, "route-3", "site-27", "site-2"),
		*NewRawRoute("route-4", enum.RouteTypeMountain, "route-4", "site-2", "site-38"),
		*NewRawRoute("route-5", enum.RouteTypeMountain, "route-5", "site-2", "site-28"),
		*NewRawRoute("route-6", enum.RouteTypeMountain, "route-6", "site-28", "site-32"),
		*NewRawRoute("route-7", enum.RouteTypeMountain, "route-7", "site-32", "site-3"),
		*NewRawRoute("route-8", enum.RouteTypeMountain, "route-8", "site-32", "site-14"),
		*NewRawRoute("route-9", enum.RouteTypeMountain, "route-9", "site-3", "site-39"),
		*NewRawRoute("route-10", enum.RouteTypeMountain, "route-10", "site-3", "site-29"),
		*NewRawRoute("route-11", enum.RouteTypeMountain, "route-11", "site-3", "site-33"),
		*NewRawRoute("route-12", enum.RouteTypePlain, "route-12", "site-29", "site-16"),
		*NewRawRoute("route-13", enum.RouteTypeMountain, "route-13", "site-29", "site-34"),
		*NewRawRoute("route-14", enum.RouteTypePlain, "route-14", "site-15", "site-5"),
		*NewRawRoute("route-15", enum.RouteTypePlain, "route-15", "site-33", "site-5"),
		*NewRawRoute("route-16", enum.RouteTypePlain, "route-16", "site-14", "site-15"),
		*NewRawRoute("route-17", enum.RouteTypePlain, "route-17", "site-14", "site-33"),
		*NewRawRoute("route-18", enum.RouteTypeSea, "route-18", "site-14", "site-13"),
		*NewRawRoute("route-19", enum.RouteTypePlain, "route-19", "site-16", "site-34"),
		*NewRawRoute("route-20", enum.RouteTypePlain, "route-20", "site-16", "site-4"),
		*NewRawRoute("route-21", enum.RouteTypePlain, "route-21", "site-16", "site-17"),
		*NewRawRoute("route-22", enum.RouteTypePlain, "route-22", "site-16", "site-5"),
		*NewRawRoute("route-23", enum.RouteTypePlain, "route-23", "site-34", "site-10"),
		*NewRawRoute("route-24", enum.RouteTypePlain, "route-24", "site-34", "site-18"),
		*NewRawRoute("route-25", enum.RouteTypePlain, "route-25", "site-4", "site-18"),
		*NewRawRoute("route-26", enum.RouteTypePlain, "route-26", "site-4", "site-17"),
		*NewRawRoute("route-27", enum.RouteTypePlain, "route-27", "site-17", "site-15"),
		*NewRawRoute("route-28", enum.RouteTypePlain, "route-28", "site-17", "site-19"),
		*NewRawRoute("route-29", enum.RouteTypePlain, "route-29", "site-15", "site-19"),
		*NewRawRoute("route-30", enum.RouteTypePlain, "route-30", "site-19", "site-18"),
		*NewRawRoute("route-31", enum.RouteTypeSea, "route-31", "site-10", "site-9"),
		*NewRawRoute("route-32", enum.RouteTypeSea, "route-32", "site-10", "site-31"),
		*NewRawRoute("route-33", enum.RouteTypePlain, "route-33", "site-9", "site-1"),
		*NewRawRoute("route-34", enum.RouteTypeSea, "route-34", "site-9", "site-31"),
		*NewRawRoute("route-35", enum.RouteTypePlain, "route-35", "site-1", "site-8"),
		*NewRawRoute("route-36", enum.RouteTypeSea, "route-36", "site-31", "site-11"),
		*NewRawRoute("route-37", enum.RouteTypeSea, "route-37", "site-8", "site-7"),
		*NewRawRoute("route-38", enum.RouteTypeSea, "route-38", "site-18", "site-30"),
		*NewRawRoute("route-39", enum.RouteTypeSea, "route-39", "site-30", "site-20"),
		*NewRawRoute("route-40", enum.RouteTypeSea, "route-40", "site-20", "site-21"),
		*NewRawRoute("route-41", enum.RouteTypeSea, "route-41", "site-20", "site-22"),
		*NewRawRoute("route-42", enum.RouteTypeSea, "route-42", "site-20", "site-35"),
		*NewRawRoute("route-43", enum.RouteTypeSea, "route-43", "site-21", "site-22"),
		*NewRawRoute("route-44", enum.RouteTypeSea, "route-44", "site-35", "site-22"),
		*NewRawRoute("route-45", enum.RouteTypeSea, "route-45", "site-22", "site-6"),
		*NewRawRoute("route-46", enum.RouteTypeSea, "route-46", "site-22", "site-25"),
		*NewRawRoute("route-47", enum.RouteTypePlain, "route-47", "site-6", "site-36"),
		*NewRawRoute("route-48", enum.RouteTypeSea, "route-48", "site-25", "site-26"),
		*NewRawRoute("route-49", enum.RouteTypeSea, "route-49", "site-25", "site-41"),
		*NewRawRoute("route-50", enum.RouteTypePlain, "route-50", "site-36", "site-23"),
		*NewRawRoute("route-51", enum.RouteTypePlain, "route-51", "site-36", "site-40"),
		*NewRawRoute("route-52", enum.RouteTypeSea, "route-52", "site-23", "site-37"),
		*NewRawRoute("route-53", enum.RouteTypeSea, "route-53", "site-37", "site-24"),
	}
}
