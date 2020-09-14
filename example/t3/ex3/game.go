package ex3

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type (
	// Game ...
	Game struct {
		bg           Scene
		currentScene Scene
		WindowSize   Size
		myUnit       Unit
		units        []Unit
	}
)

// NewGame ...
func NewGame(w, h int) (*Game, error) {
	size := Size{
		Width:  w,
		Height: h,
	}

	backGround, _ := NewBackGround(size)
	g := &Game{
		bg:         backGround,
		WindowSize: size,
	}

	// 初期化時のシーンを登録
	// 仮のユニットを作成
	s := NewBattleScene(size)
	g.currentScene = s

	// Unit
	u, _ := NewMyUnit(g)
	g.myUnit = u

	u2, _ := NewUnit(g)
	g.units = append(g.units, u2)

	return g, nil
}

// Update ...
func (g *Game) Update(r *ebiten.Image) error {
	const d = 16

	sw, sh := r.Size()
	dbg := fmt.Sprintf("screen size: %d, %d", sw, sh)

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
			u, _ := NewDebris(0, g)
			// 生成オブジェクトの衝突判定
			col := false
			for _, unit := range g.units {
				if CollisionUnit(u, unit) {
					col = true
					break
				}
			}
			if !col {
				g.units = append(g.units, u)
			}
		}
	}

	if err := g.bg.Update(); err != nil {
		return err
	}
	if err := g.currentScene.Update(); err != nil {
		return err
	}
	if err := g.myUnit.Update(); err != nil {
		return err
	}
	for _, u := range g.units {
		if err := u.Update(); err != nil {
			return err
		}
	}
	// ユニットの衝突判定
	for _, u := range g.units {
		if CollisionUnit(g.myUnit, u) {
			g.myUnit.Collision(&u)
			u.Collision(&g.myUnit)
		}
		// _ = Dot(g.myUnit, u)
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	g.bg.Draw(r)
	g.currentScene.Draw(r)
	g.myUnit.Draw(r)
	for _, u := range g.units {
		u.Draw(r)
	}

	ebitenutil.DebugPrint(r, dbg)
	return nil
}

// CollisionUnit unit同士の衝突状態を返す
func CollisionUnit(unit1, unit2 Unit) bool {
	u1, u2 := unit1.Circle(), unit2.Circle()
	// (xc1-xc2)^2 + (yc1-yc2)^2 ≦ (r1+r2)^2
	var dx, dy, dr float64 = float64(u1.x - u2.x), float64(u1.y - u2.y), float64(u1.r + u2.r)
	if (dx*dx + dy*dy) <= dr*dr {
		return true
	}
	return false
}

// Dot ...
func Dot(unit1, unit2 Unit) float64 {
	//	x1*x2 + y1*y2
	p := unit1.Circle().x*unit2.Circle().x + unit1.Circle().y*unit2.Circle().y
	// log.Printf("Dot=%f", p)
	return p
}
