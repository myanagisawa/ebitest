package ex3

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/myanagisawa/ebitest/utils"
)

const (
	d = 16
)

type (
	// Game ...
	Game struct {
		currentScene Scene
		WindowSize   Size
	}
)

var (
	fface10White *LabelFace
	fface10Red   *LabelFace
	fface10Green *LabelFace
	images       map[string]image.Image

	paused = false
)

func init() {
	fface10White = NewLabelFace(10, color.White)
	fface10Red = NewLabelFace(10, color.RGBA{255, 0, 0, 255})
	fface10Green = NewLabelFace(10, color.RGBA{0, 255, 0, 255})

	images = make(map[string]image.Image)
	img, _ := utils.OrientationImage("resources/system_images/unit-0.png")
	images["unit-0"] = img

	img, _ = utils.OrientationImage("resources/system_images/unit-1.png")
	images["unit-1"] = img

	img, _ = utils.OrientationImage("resources/system_images/unit-2.png")
	images["unit-2"] = img

	img, _ = utils.OrientationImage("resources/system_images/unit-del.png")
	images["unit-del"] = img

	img, _ = utils.OrientationImage("resources/system_images/search-1.png")
	images["search-1"] = img

	img, _ = utils.OrientationImage("resources/system_images/pin.png")
	images["pin"] = img
}

// NewGame ...
func NewGame(w, h int) (*Game, error) {
	rand.Seed(time.Now().UnixNano()) //Seed

	size := Size{
		Width:  w,
		Height: h,
	}

	g := &Game{
		WindowSize: size,
	}

	// 初期化時のシーンを登録
	// 仮のユニットを作成
	s := NewBattleScene(size)

	// teamを追加
	t1 := s.AddTeam(&Point{X: 200, Y: 400}, true)
	t2 := s.AddTeam(&Point{X: 800, Y: 300}, false)
	t3 := s.AddTeam(&Point{X: 900, Y: 600}, false)

	// 対決状態を設定
	t1.Enemies = []*Team{t2, t3}
	t2.Enemies = []*Team{t1}
	t3.Enemies = []*Team{t1}

	// 同盟状態を設定
	t2.Alliances = []*Team{t3}
	t3.Alliances = []*Team{t2}

	// Unitを作成
	addUnitToTeam(t1, rand.Intn(7)+5)
	addUnitToTeam(t2, rand.Intn(7)+3)
	addUnitToTeam(t3, rand.Intn(7)+3)

	s.InitWindows()

	g.currentScene = s

	return g, nil
}

func addUnitToTeam(t *Team, num int) {
	rand.Seed(time.Now().UnixNano()) //Seed

	enemy := t.Enemies[0]
	// unitの向きを敵チームに向ける
	x1, y1 := t.Location.X, t.Location.Y
	x2, y2 := enemy.Location.X, enemy.Location.Y
	dx, dy := x2-x1, -(y2 - y1) // 画面の上側をY座標＋とするので、Y座標は符号を入れ替える
	// radianを取得
	n := math.Atan2(float64(dy), float64(dx))
	// radian ->degreeに変換
	d := n * 180 / math.Pi
	a := int(d)

	for i := 0; i < num; i++ {
		l := fmt.Sprintf("team_%d_%d", t.No, i)
		size := 10
		if i == 0 {
			size = 15
		}
		x, y := t.Location.X+(size*2*i), t.Location.Y+(size*2*i)
		s := 1
		hp := rand.Intn(20) + 5
		u, err := NewUnit(t.Parent, t.No, hp, size, l, x, y, a, s, 100)
		if err != nil {
			panic(err)
		}
		t.Units = append(t.Units, u)
	}
	// log.Printf("team: %#v", t)
}

// Update ...
func (g *Game) Update(r *ebiten.Image) error {
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

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		fmt.Println("Game::C")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		fmt.Println("Game::P")
		paused = !paused
	}

	if !paused {
		if err := g.currentScene.Update(); err != nil {
			return err
		}
	}

	g.currentScene.Draw(r)

	if paused {
		dbg += "\nPause"
	}
	ebitenutil.DebugPrint(r, dbg)
	return nil
}
