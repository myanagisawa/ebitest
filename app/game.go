package app

import (
	"github.com/myanagisawa/ebitest/app/g"
	"github.com/myanagisawa/ebitest/app/menu"
	"github.com/myanagisawa/ebitest/app/obj"
	"github.com/myanagisawa/ebitest/app/wmap"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/game"
)

// NewGameManager ...
func NewGameManager(screenWidth, screenHeight int) *game.Manager {
	gm := game.NewManager(screenWidth, screenHeight)

	// Dataロード
	sites := obj.CreateSites(g.Master.Sites)
	routes := obj.CreateRoutes(g.Master.Routes, sites)

	gm.AddDataSet(enum.DataTypeSite, sites)
	gm.AddDataSet(enum.DataTypeRoute, routes)

	// メニュー画面ロード
	menu := menu.NewScene(gm)
	gm.SetScene(enum.MenuEnum, menu)
	// MAP画面ロード
	wm := wmap.NewScene(gm)
	gm.SetScene(enum.MapEnum, wm)

	// 起動画面
	gm.TransitionTo(enum.MenuEnum)
	return gm
}
