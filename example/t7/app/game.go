package app

import (
	"log"

	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/app/m"
	"github.com/myanagisawa/ebitest/example/t7/app/menu"
	"github.com/myanagisawa/ebitest/example/t7/app/wmap"
	"github.com/myanagisawa/ebitest/example/t7/lib/game"
)

var (
	// rawData マスタデータ
	rawData *m.MasterData
)

// NewGameManager ...
func NewGameManager(screenWidth, screenHeight int) *game.Manager {
	// 定義データロード
	log.Printf("MASTERデータ読み込み...")
	rawData = m.LoadMaster()
	log.Printf("MASTERデータ読み込み...done")

	// ゲームデータ初期化
	log.Printf("ゲームデータ初期化...")
	sites := m.CreateSites(rawData.RawSites)
	routes := m.CreateRoutes(rawData.RawRoutes, sites)
	log.Printf("ゲームデータ初期化...done")

	gm := game.NewManager(screenWidth, screenHeight)

	// ゲームデータ登録
	gm.SetGameData(enum.DataTypeSite, sites)
	gm.SetGameData(enum.DataTypeRoute, routes)

	// メニュー画面ロード
	menu := menu.NewScene(gm)
	gm.SetScene(enum.MenuEnum, menu)
	// マップ画面ロード
	wmap := wmap.NewScene(gm)
	gm.SetScene(enum.MapEnum, wmap)
	// 起動画面
	gm.TransitionTo(enum.MenuEnum)
	return gm
}
