package kitchen

import (
	"fmt"
	"image/color"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
)

type (
	// Game ...
	Game struct {
		currentScene Scene
		WindowSize   Size
	}

	// Size ...
	Size struct {
		Width  int
		Height int
	}
)

var (
	uiFont        font.Face
	uiFontMHeight int
)

// NewGame ...
func NewGame(w, h int) (*Game, error) {

	g := &Game{
		WindowSize: Size{
			Width:  w,
			Height: h,
		},
	}

	// 初期化時のシーンを登録
	sink, _ := NewSink(&g.WindowSize)
	g.currentScene = sink

	// ebitenフォントのテスト
	tt, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	uiFont = truetype.NewFace(tt, &truetype.Options{
		Size:    48,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	b, _, _ := uiFont.GlyphBounds('M')
	uiFontMHeight = (b.Max.Y - b.Min.Y).Ceil()

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

	r.Fill(color.RGBA{0x90, 0x7e, 0xb4, 0xdd})

	str := fmt.Sprintf("w=%d, h=%d", sw, sh)
	text.Draw(r, str, uiFont, 50, 50, color.White)

	if err := g.currentScene.Update(); err != nil {
		return err
	}
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	g.currentScene.Draw(r)

	//ebitenutil.DebugPrint(r, dbg)
	return nil
}
