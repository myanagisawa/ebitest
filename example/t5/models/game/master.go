package game

import (
	"fmt"

	"github.com/myanagisawa/ebitest/example/t5/ebitest"
	"github.com/myanagisawa/ebitest/example/t5/enum"
)

// MasterData ...
type MasterData struct {
	Sites  []*MSite
	Routes []*MRoute
}

// NewMasterData ...
func NewMasterData() *MasterData {
	m := &MasterData{}
	m.Sites = DefSites()
	m.Routes = DefRoutes()
	return m
}

// MSite マップ上の拠点情報
type MSite struct {
	Code     string
	Type     enum.SiteTypeEnum
	Name     string
	Location *ebitest.Point
}

// NewMSite ...
func NewMSite(code string, t enum.SiteTypeEnum, name string, loc *ebitest.Point) *MSite {
	if code == "" || loc.X() > 1 || loc.Y() > 1 {
		panic(fmt.Sprintf("Site初期化エラー: code=%#v, loc=%#v", code, loc))
	}

	return &MSite{
		Code:     code,
		Type:     t,
		Name:     name,
		Location: loc,
	}
}

// MRoute マップ上の経路情報
type MRoute struct {
	Code  string
	Type  enum.RouteTypeEnum
	Name  string
	Site1 string
	Site2 string
}

// NewMRoute ...
func NewMRoute(code string, t enum.RouteTypeEnum, name, site1, site2 string) *MRoute {
	if code == "" || site1 == "" || site2 == "" {
		panic(fmt.Sprintf("Route初期化エラー: code=%#v, site1=%#v, site2=%#v", code, site1, site2))
	}

	return &MRoute{
		Code:  code,
		Type:  t,
		Name:  name,
		Site1: site1,
		Site2: site2,
	}
}

// DefSites ...
func DefSites() []*MSite {
	return []*MSite{
		NewMSite("site-1", enum.SiteTypeBase, "base-1", ebitest.NewPoint(0.318, 0.707)),
		NewMSite("site-2", enum.SiteTypeBase, "base-2", ebitest.NewPoint(0.161, 0.761)),
		NewMSite("site-3", enum.SiteTypeBase, "base-3", ebitest.NewPoint(0.237, 0.826)),
		NewMSite("site-4", enum.SiteTypeBase, "base-4", ebitest.NewPoint(0.345, 0.847)),
		NewMSite("site-5", enum.SiteTypeBase, "base-5", ebitest.NewPoint(0.255, 0.883)),
		NewMSite("site-6", enum.SiteTypeBase, "base-6", ebitest.NewPoint(0.573, 0.780)),
		NewMSite("site-7", enum.SiteTypeCity, "city-1", ebitest.NewPoint(0.306, 0.623)),
		NewMSite("site-8", enum.SiteTypeCity, "city-2", ebitest.NewPoint(0.302, 0.663)),
		NewMSite("site-9", enum.SiteTypeCity, "city-3", ebitest.NewPoint(0.295, 0.726)),
		NewMSite("site-10", enum.SiteTypeCity, "city-4", ebitest.NewPoint(0.276, 0.751)),
		NewMSite("site-11", enum.SiteTypeCity, "city-5", ebitest.NewPoint(0.181, 0.648)),
		NewMSite("site-12", enum.SiteTypeCity, "city-6", ebitest.NewPoint(0.181, 0.598)),
		NewMSite("site-13", enum.SiteTypeCity, "city-7", ebitest.NewPoint(0.120, 0.847)),
		NewMSite("site-14", enum.SiteTypeCity, "city-8", ebitest.NewPoint(0.193, 0.863)),
		NewMSite("site-15", enum.SiteTypeCity, "city-9", ebitest.NewPoint(0.263, 0.953)),
		NewMSite("site-16", enum.SiteTypeCity, "city-10", ebitest.NewPoint(0.279, 0.842)),
		NewMSite("site-17", enum.SiteTypeCity, "city-11", ebitest.NewPoint(0.314, 0.893)),
		NewMSite("site-18", enum.SiteTypeCity, "city-12", ebitest.NewPoint(0.374, 0.835)),
		NewMSite("site-19", enum.SiteTypeCity, "city-13", ebitest.NewPoint(0.378, 0.932)),
		NewMSite("site-20", enum.SiteTypeCity, "city-14", ebitest.NewPoint(0.448, 0.863)),
		NewMSite("site-21", enum.SiteTypeCity, "city-15", ebitest.NewPoint(0.479, 0.824)),
		NewMSite("site-22", enum.SiteTypeCity, "city-16", ebitest.NewPoint(0.563, 0.852)),
		NewMSite("site-23", enum.SiteTypeCity, "city-17", ebitest.NewPoint(0.565, 0.708)),
		NewMSite("site-24", enum.SiteTypeCity, "city-18", ebitest.NewPoint(0.600, 0.613)),
		NewMSite("site-25", enum.SiteTypeCity, "city-19", ebitest.NewPoint(0.652, 0.856)),
		NewMSite("site-26", enum.SiteTypeCity, "city-20", ebitest.NewPoint(0.698, 0.805)),
		NewMSite("site-27", enum.SiteTypeFort, "fort-1", ebitest.NewPoint(0.181, 0.718)),
		NewMSite("site-28", enum.SiteTypeFort, "fort-2", ebitest.NewPoint(0.189, 0.785)),
		NewMSite("site-29", enum.SiteTypeFort, "fort-3", ebitest.NewPoint(0.252, 0.801)),
		NewMSite("site-30", enum.SiteTypeFort, "fort-4", ebitest.NewPoint(0.413, 0.834)),
		NewMSite("site-31", enum.SiteTypePoint, "point-1", ebitest.NewPoint(0.245, 0.714)),
		NewMSite("site-32", enum.SiteTypePoint, "point-2", ebitest.NewPoint(0.209, 0.822)),
		NewMSite("site-33", enum.SiteTypePoint, "point-3", ebitest.NewPoint(0.230, 0.868)),
		NewMSite("site-34", enum.SiteTypePoint, "point-4", ebitest.NewPoint(0.302, 0.807)),
		NewMSite("site-35", enum.SiteTypePoint, "point-5", ebitest.NewPoint(0.484, 0.918)),
		NewMSite("site-36", enum.SiteTypePoint, "point-6", ebitest.NewPoint(0.566, 0.751)),
		NewMSite("site-37", enum.SiteTypePoint, "point-7", ebitest.NewPoint(0.580, 0.656)),
		NewMSite("site-38", enum.SiteTypeRemains, "remains-1", ebitest.NewPoint(0.127, 0.730)),
		NewMSite("site-39", enum.SiteTypeRemains, "remains-2", ebitest.NewPoint(0.223, 0.786)),
		NewMSite("site-40", enum.SiteTypeRemains, "remains-3", ebitest.NewPoint(0.580, 0.747)),
		NewMSite("site-41", enum.SiteTypeRemains, "remains-4", ebitest.NewPoint(0.636, 0.793)),
	}
}

// DefRoutes ...
func DefRoutes() []*MRoute {
	return []*MRoute{
		NewMRoute("route-1", enum.RouteTypeSea, "route-1", "site-11", "site-12"),
		NewMRoute("route-2", enum.RouteTypePlain, "route-2", "site-12", "site-27"),
		NewMRoute("route-3", enum.RouteTypeMountain, "route-3", "site-27", "site-2"),
		NewMRoute("route-4", enum.RouteTypeMountain, "route-4", "site-2", "site-38"),
		NewMRoute("route-5", enum.RouteTypeMountain, "route-5", "site-2", "site-28"),
		NewMRoute("route-6", enum.RouteTypeMountain, "route-6", "site-28", "site-32"),
		NewMRoute("route-7", enum.RouteTypeMountain, "route-7", "site-32", "site-3"),
		NewMRoute("route-8", enum.RouteTypeMountain, "route-8", "site-32", "site-14"),
		NewMRoute("route-9", enum.RouteTypeMountain, "route-9", "site-3", "site-39"),
		NewMRoute("route-10", enum.RouteTypeMountain, "route-10", "site-3", "site-29"),
		NewMRoute("route-11", enum.RouteTypeMountain, "route-11", "site-3", "site-33"),
		NewMRoute("route-12", enum.RouteTypePlain, "route-12", "site-29", "site-16"),
		NewMRoute("route-13", enum.RouteTypeMountain, "route-13", "site-29", "site-34"),
		NewMRoute("route-14", enum.RouteTypePlain, "route-14", "site-33", "site-14"),
		NewMRoute("route-15", enum.RouteTypePlain, "route-15", "site-33", "site-5"),
		NewMRoute("route-16", enum.RouteTypePlain, "route-16", "site-14", "site-15"),
		NewMRoute("route-17", enum.RouteTypePlain, "route-17", "site-14", "site-33"),
		NewMRoute("route-18", enum.RouteTypeSea, "route-18", "site-14", "site-13"),
		NewMRoute("route-19", enum.RouteTypePlain, "route-19", "site-16", "site-34"),
		NewMRoute("route-20", enum.RouteTypePlain, "route-20", "site-16", "site-4"),
		NewMRoute("route-21", enum.RouteTypePlain, "route-21", "site-16", "site-17"),
		NewMRoute("route-22", enum.RouteTypePlain, "route-22", "site-16", "site-5"),
		NewMRoute("route-23", enum.RouteTypePlain, "route-23", "site-34", "site-10"),
		NewMRoute("route-24", enum.RouteTypePlain, "route-24", "site-34", "site-18"),
		NewMRoute("route-25", enum.RouteTypePlain, "route-25", "site-4", "site-18"),
		NewMRoute("route-26", enum.RouteTypePlain, "route-26", "site-4", "site-17"),
		NewMRoute("route-27", enum.RouteTypePlain, "route-27", "site-17", "site-15"),
		NewMRoute("route-28", enum.RouteTypePlain, "route-28", "site-17", "site-19"),
		NewMRoute("route-29", enum.RouteTypePlain, "route-29", "site-15", "site-19"),
		NewMRoute("route-30", enum.RouteTypePlain, "route-30", "site-19", "site-18"),
		NewMRoute("route-31", enum.RouteTypeSea, "route-31", "site-10", "site-9"),
		NewMRoute("route-32", enum.RouteTypeSea, "route-32", "site-10", "site-31"),
		NewMRoute("route-33", enum.RouteTypePlain, "route-33", "site-9", "site-1"),
		NewMRoute("route-34", enum.RouteTypeSea, "route-34", "site-9", "site-31"),
		NewMRoute("route-35", enum.RouteTypePlain, "route-35", "site-1", "site-8"),
		NewMRoute("route-36", enum.RouteTypeSea, "route-36", "site-31", "site-12"),
		NewMRoute("route-37", enum.RouteTypeSea, "route-37", "site-8", "site-7"),
		NewMRoute("route-38", enum.RouteTypeSea, "route-38", "site-18", "site-30"),
		NewMRoute("route-39", enum.RouteTypeSea, "route-39", "site-30", "site-20"),
		NewMRoute("route-40", enum.RouteTypeSea, "route-40", "site-20", "site-21"),
		NewMRoute("route-41", enum.RouteTypeSea, "route-41", "site-20", "site-22"),
		NewMRoute("route-42", enum.RouteTypeSea, "route-42", "site-20", "site-35"),
		NewMRoute("route-43", enum.RouteTypeSea, "route-43", "site-21", "site-22"),
		NewMRoute("route-44", enum.RouteTypeSea, "route-44", "site-35", "site-22"),
		NewMRoute("route-45", enum.RouteTypeSea, "route-45", "site-22", "site-6"),
		NewMRoute("route-46", enum.RouteTypeSea, "route-46", "site-22", "site-25"),
		NewMRoute("route-47", enum.RouteTypePlain, "route-47", "site-6", "site-36"),
		NewMRoute("route-48", enum.RouteTypeSea, "route-48", "site-25", "site-26"),
		NewMRoute("route-49", enum.RouteTypeSea, "route-49", "site-25", "site-41"),
		NewMRoute("route-50", enum.RouteTypePlain, "route-50", "site-36", "site-23"),
		NewMRoute("route-51", enum.RouteTypePlain, "route-51", "site-36", "site-40"),
		NewMRoute("route-52", enum.RouteTypeSea, "route-52", "site-23", "site-37"),
		NewMRoute("route-53", enum.RouteTypeSea, "route-53", "site-37", "site-24"),
	}
}
