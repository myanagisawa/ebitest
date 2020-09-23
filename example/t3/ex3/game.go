package ex3

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/myanagisawa/ebitest/utils"
)

type (
	// Game ...
	Game struct {
		bg           Scene
		currentScene Scene
		WindowSize   Size
		teams        []Team
		myUnit       Unit
		units        []Unit
	}

	// Team ...
	Team struct {
		No        int
		Units     []Unit
		Enemies   []*Team
		Alliances []*Team
		Location  Point
	}
)

var (
	fface10White *LabelFace
	fface10Red   *LabelFace
	images       map[string]*ebiten.Image
)

func init() {
	fface10White = NewLabelFace(10, color.White)
	fface10Red = NewLabelFace(10, color.RGBA{255, 0, 0, 255})

	images = make(map[string]*ebiten.Image)
	img, _ := utils.OrientationImage("resources/system_images/unit-0.png")
	eimg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	images["unit-0"] = eimg

	img, _ = utils.OrientationImage("resources/system_images/unit-1.png")
	eimg, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	images["unit-1"] = eimg

	img, _ = utils.OrientationImage("resources/system_images/unit-2.png")
	eimg, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	images["unit-2"] = eimg

	img, _ = utils.OrientationImage("resources/system_images/unit-del.png")
	eimg, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	images["unit-del"] = eimg

	img, _ = utils.OrientationImage("resources/system_images/search-1.png")
	eimg, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	images["search-1"] = eimg

}

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

	// 自チーム、敵チーム1、敵チーム2を作成
	teams := make([]Team, 3)
	teams[0] = Team{No: 0, Location: Point{X: 200, Y: 400}}
	teams[1] = Team{No: 1, Location: Point{X: 800, Y: 300}}
	teams[2] = Team{No: 2, Location: Point{X: 600, Y: 600}}

	// 対決状態を設定
	teams[0].Enemies = []*Team{&teams[1], &teams[2]}
	teams[1].Enemies = []*Team{&teams[0]}
	teams[2].Enemies = []*Team{&teams[0]}

	// 同盟状態を設定
	teams[1].Alliances = []*Team{&teams[2]}
	teams[2].Alliances = []*Team{&teams[1]}

	// Unitを作成
	for t, team := range teams {
		enemy := team.Enemies[0]
		// unitの向きを敵チームに向ける
		x1, y1 := team.Location.X, team.Location.Y
		x2, y2 := enemy.Location.X, enemy.Location.Y
		dx, dy := x2-x1, -(y2 - y1) // 画面の上側をY座標＋とするので、Y座標は符号を入れ替える
		// radianを取得
		n := math.Atan2(float64(dy), float64(dx))
		// radian ->degreeに変換
		d := n * 180 / math.Pi
		a := int(d)

		for i := 0; i < 5; i++ {
			l := fmt.Sprintf("team_%d_%d", team.No, i)
			size := 10
			x, y := team.Location.X+(size*2*i), team.Location.Y+(size*2*i)
			s := 1
			u, err := NewUnit(g, team.No, 5, size, l, x, y, a, s, 100)
			if err != nil {
				panic(err)
			}
			team.Units = append(team.Units, u)
		}
		teams[t] = team
		log.Printf("team: %#v", team)
	}

	g.teams = teams

	// // Unit
	// u, _ := NewMyUnit(g)
	// g.myUnit = u

	// u2, _ := NewUnit(g)
	// g.units = append(g.units, u2)

	return g, nil
}

// Update ...
func (g *Game) Update(r *ebiten.Image) error {
	const d = 16

	sw, sh := r.Size()
	dbg := fmt.Sprintf("screen size: %d, %d\nFPS: %0.2f", sw, sh, ebiten.CurrentFPS())

	// 停止、アクティブ実装
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
	}

	if err := g.bg.Update(); err != nil {
		return err
	}

	for _, team := range g.teams {
		for _, u := range team.Units {
			if u.GetStatus() != 0 {
				continue
			}
			// ユニットのレーダー捕捉判定
			u.SetCaptured(nil)
			for _, et := range team.Enemies {
				// log.Printf("et: %d, %v", et.No, et.Units)
				captureUnit(u, et.Units)
			}

			// ユニットの衝突判定
			for _, et := range team.Enemies {
				for _, eu := range et.Units {
					if eu.GetStatus() != 0 {
						continue
					}
					// log.Printf("Collision")
					if CollisionUnit(u, eu) {
						u.Collision(&eu)
					}
				}
			}

			if err := u.Update(); err != nil {
				return err
			}
		}
	}

	if err := g.currentScene.Update(); err != nil {
		return err
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	g.bg.Draw(r)

	for _, team := range g.teams {
		for _, u := range team.Units {
			u.Draw(r)
		}
	}

	g.currentScene.Draw(r)

	ebitenutil.DebugPrint(r, dbg)
	return nil
}

// CollisionUnit unit同士の衝突状態を返す
func CollisionUnit(unit1, unit2 Unit) bool {
	x1, y1 := unit1.GetCenter()
	x2, y2 := unit2.GetCenter()
	e1, e2 := unit1.GetEntity(), unit2.GetEntity()
	// (xc1-xc2)^2 + (yc1-yc2)^2 ≦ (r1+r2)^2
	var dx, dy, dr float64 = float64(x1 - x2), float64(y1 - y2), float64(e1.R() + e2.R())
	if (dx*dx + dy*dy) <= dr*dr {
		return true
	}
	return false
}

// captureUnit unitの索敵範囲に入ったunitsを取得する
func captureUnit(unit Unit, units []Unit) {
	x1, y1 := unit.GetCenter()
	c1 := unit.GetRader()
	captured := []Unit{}
	for _, u := range units {
		if u.GetStatus() != 0 {
			continue
		}
		x2, y2 := u.GetCenter()
		e := u.GetEntity()

		var dx, dy, dr float64 = float64(x1 - x2), float64(y1 - y2), float64(c1.R() + e.R())
		if (dx*dx + dy*dy) <= dr*dr {
			// レーダー捕捉
			captured = append(captured, u)
			// log.Printf("captured!!")
		}
	}
	unit.SetCaptured(captured)
}

// Dot ...
func Dot(unit1, unit2 Unit) float64 {
	x1, y1 := unit1.GetCenter()
	x2, y2 := unit2.GetCenter()
	//	x1*x2 + y1*y2
	p := x1*x2 + y1*y2
	// log.Printf("Dot=%f", p)
	return p
}
