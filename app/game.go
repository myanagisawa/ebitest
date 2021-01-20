package app

import (
	"github.com/myanagisawa/ebitest/app/menu"
	"github.com/myanagisawa/ebitest/app/wmap"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/game"
)

// NewGameManager ...
func NewGameManager(screenWidth, screenHeight int) *game.Manager {
	gm := game.NewManager(screenWidth, screenHeight)
	// メニュー画面ロード
	gm.SetScene(enum.MenuEnum, menu.NewScene(gm))
	// MAP画面ロード
	gm.SetScene(enum.MapEnum, wmap.NewScene(gm))

	// 起動画面
	gm.TransitionTo(enum.MenuEnum)
	return gm
}
