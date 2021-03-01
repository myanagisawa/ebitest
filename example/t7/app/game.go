package app

import (
	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/app/menu"
	"github.com/myanagisawa/ebitest/example/t7/lib/game"
)

// NewGameManager ...
func NewGameManager(screenWidth, screenHeight int) *game.Manager {
	gm := game.NewManager(screenWidth, screenHeight)

	// メニュー画面ロード
	menu := menu.NewScene()
	gm.SetScene(enum.MenuEnum, menu)
	// 起動画面
	gm.TransitionTo(enum.MenuEnum)
	return gm
}
