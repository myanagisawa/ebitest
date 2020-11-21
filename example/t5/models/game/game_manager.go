package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/example/t5/ebitest"
	"github.com/myanagisawa/ebitest/example/t5/enum"
	"github.com/myanagisawa/ebitest/example/t5/interfaces"
	"github.com/myanagisawa/ebitest/example/t5/models/scene"
	"github.com/myanagisawa/ebitest/example/t5/parts"
)

/*
	- マスタからゲームデータを生成する生成処理: OK
	- レイヤサイズを画面サイズに合わせるオプション
	- ゲームデータの画面上への描画
*/

type (
	// Manager ...
	Manager struct {
		master       *MasterData
		currentScene interfaces.Scene
		scenes       map[enum.SceneEnum]interfaces.Scene
	}
)

func getSiteByCode(code string, sites []*parts.Site) *parts.Site {
	for _, site := range sites {
		if site.Code == code {
			return site
		}
	}
	return nil
}

// GetSites master.MSiteからSiteを作成します
func (g *Manager) GetSites() []*parts.Site {
	sites := make([]*parts.Site, len(g.master.Sites))
	for i, site := range g.master.Sites {
		sites[i] = parts.NewSite(site.Code, site.Type, site.Name, site.Location)
	}
	return sites
}

// GetRoutes master.MSiteからSiteを作成します
func (g *Manager) GetRoutes(sites []*parts.Site) []*parts.Route {
	routes := make([]*parts.Route, len(g.master.Routes))
	for i, route := range g.master.Routes {
		routes[i] = parts.NewRoute(
			route.Code,
			route.Type,
			route.Name,
			getSiteByCode(route.Site1, sites),
			getSiteByCode(route.Site2, sites))
	}
	return routes
}

// NewManager ...
func NewManager(screenWidth, screenHeight int) *Manager {
	ebitest.Width, ebitest.Height = screenWidth, screenHeight

	gm := &Manager{
		master: NewMasterData(),
	}

	scenes := map[enum.SceneEnum]interfaces.Scene{}
	scenes[enum.MainMenuEnum] = scene.NewMainMenu(gm)
	scenes[enum.MapEnum] = scene.NewMap(gm)
	gm.scenes = scenes

	// MainMenuを表示
	gm.TransitionTo(enum.MapEnum)
	return gm
}

// TransitionTo ...
func (g *Manager) TransitionTo(t enum.SceneEnum) {
	// var s interfaces.Scene
	// switch t {
	// case enum.MainMenuEnum:
	// 	s = scene.NewMainMenu(g)
	// case enum.MapEnum:
	// 	s = scene.NewMap(g)
	// default:
	// 	panic(fmt.Sprintf("invalid SceneEnum: %d", t))
	// }
	g.currentScene = g.scenes[t]
}

// SetCurrentScene ...
func (g *Manager) SetCurrentScene(s interfaces.Scene) {
	g.currentScene = s
}

// Update ...
func (g *Manager) Update() error {
	if g.currentScene != nil {
		return g.currentScene.Update()
	}
	return nil
}

// Draw ...
func (g *Manager) Draw(screen *ebiten.Image) {
	if g.currentScene != nil {
		g.currentScene.Draw(screen)
	}
}

// Layout ...
func (g *Manager) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ebitest.Width, ebitest.Height
}
