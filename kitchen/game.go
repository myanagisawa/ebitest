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
		myCoin       Coin
		coins        []Coin
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
	c, _ := NewMyCoin()
	g.myCoin = c

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
		for i := 0; i < 10; i++ {
			c, _ := NewDebris(0)
			g.coins = append(g.coins, c)
		}
	}

	if err := g.bg.Update(); err != nil {
		return err
	}
	if err := g.currentScene.Update(); err != nil {
		return err
	}
	if err := g.debugText.Update(); err != nil {
		return err
	}
	if err := g.myCoin.Update(); err != nil {
		return err
	}
	for _, c := range g.coins {
		if err := c.Update(); err != nil {
			return err
		}
	}
	// コインの衝突判定
	for _, c := range g.coins {
		CollisionCoin(g.myCoin, c)
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	g.bg.Draw(r)
	g.currentScene.Draw(r)
	g.debugText.Draw(r)
	g.myCoin.Draw(r)
	for _, c := range g.coins {
		c.Draw(r)
	}

	//ebitenutil.DebugPrint(r, dbg)
	return nil
}

// CollisionCoin ...
func CollisionCoin(coin1, coin2 Coin) {
	c1, c2 := coin1.Circle(), coin2.Circle()
	// (xc1-xc2)^2 + (yc1-yc2)^2 ≦ (r1+r2)^2
	var dx, dy, dr float64 = float64(c1.x - c2.x), float64(c1.y - c2.y), float64(c1.r + c2.r)
	if (dx*dx + dy*dy) <= dr*dr {
		coin1.Collision(&coin2)
		coin2.Collision(&coin1)
	}
}
