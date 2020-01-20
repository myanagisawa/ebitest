package kitchen

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type (
	// Game ...
	Game struct {
		bg           Scene
		currentScene Scene
		debugText    *DebugText
		coins        []Coin
		light        Spotlight
		WindowSize   Size
	}

	// Size ...
	Size struct {
		Width  int
		Height int
	}
)

// NewGame ...
func NewGame(w, h int) (*Game, error) {

	backGround, _ := NewBackGround()
	debugText, _ := NewDebugText()
	g := &Game{
		bg:        backGround,
		debugText: debugText,
		WindowSize: Size{
			Width:  w,
			Height: h,
		},
	}

	// 初期化時のシーンを登録
	// sink, _ := NewSink(&g.WindowSize)
	// g.currentScene = sink
	s := NewSceneImpl()
	g.currentScene = s

	// Coin
	// c, _ := NewCoin()
	// g.coin = c

	// l, _ := NewSpotlight(300.0, 300.0, 200.0, 1)
	// g.light = *l
	return g, nil
}

// Update ...
func (g *Game) Update(r *ebiten.Image) error {
	const d = 16

	sw, sh := r.Size()
	//dbg := fmt.Sprintf("screen size: %d, %d", sw, sh)

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		fmt.Println("Game::Up")
		sh += d
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		fmt.Println("Game::Down")
		if 16 < sh && d < sh {
			sh -= d
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		fmt.Println("Game::Left")
		if 16 < sw && d < sw {
			sw -= d
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		fmt.Println("Game::Right")
		sw += d
	}
	ebiten.SetScreenSize(sw, sh)

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		fmt.Println("Game::C")
		c, _ := NewCoin()
		g.coins = append(g.coins, c)
	}
	// str := fmt.Sprintf("w=%d, h=%d\n", sw, sh)
	// g.debugText.Append(str)
	// angle, speed, vx, vy, x, y := g.coin.Info()
	// str := fmt.Sprintf("coin: a=%d, s=%d, vx=%d, vy=%d, x=%d, y=%d", angle, speed, vx, vy, x, y)
	// g.debugText.Append(str)

	if err := g.bg.Update(); err != nil {
		return err
	}
	if err := g.currentScene.Update(); err != nil {
		return err
	}
	if err := g.debugText.Update(); err != nil {
		return err
	}
	for _, c := range g.coins {
		if err := c.Update(); err != nil {
			return err
		}
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	g.bg.Draw(r)
	g.currentScene.Draw(r)
	g.debugText.Draw(r)
	for _, c := range g.coins {
		c.Draw(r)
	}

	//ebitenutil.DebugPrint(r, dbg)
	return nil
}
