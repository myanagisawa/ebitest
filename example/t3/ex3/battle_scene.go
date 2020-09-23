package ex3

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

type (
	// BattleScene ...
	BattleScene struct {
		size Size
	}
)

var (
	posX, posY = 5.0, 35.0
	winWidth   = 400.0
	margin     = 3.0
	rowWidth   = winWidth - (margin * 2)
	rowHeight  = 30.0

	rownum    = 10
	winHeight = (rowHeight * float64(rownum)) + (margin * float64(rownum+1))

	winImg *ebiten.Image
	rowImg *ebiten.Image

	textAdjust int
)

// NewBattleScene ...
func NewBattleScene(s Size) *BattleScene {
	winImg = createRectImage(int(winWidth), int(winHeight), color.RGBA{0, 0, 0, 64})
	rowImg = createRectImage(int(rowWidth), int(rowHeight), color.RGBA{50, 50, 50, 64})

	textAdjust = int(rowHeight / 2)
	scene := &BattleScene{
		size: s,
	}
	return scene
}

// Update ...
func (s *BattleScene) Update() error {

	//log.Printf("BattleScene.Update")

	return nil
}

// Draw ...
func (s *BattleScene) Draw(r *ebiten.Image) {
	for i := 0; i < rownum; i++ {
		n := float64(i)
		y := (margin * (n + 1)) + (rowHeight * n)

		text.Draw(winImg, fmt.Sprintf("%d", i), fface10White.uiFont, int(margin), int(y)+textAdjust, fface10White.uiFontColor)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(posX, posY)

	r.DrawImage(winImg, op)

	// ebitenutil.DrawRect(r, posX, posY, winWidth, winHeight, color.RGBA{0, 0, 0, 64})

	// for i := 0; i < 10; i++ {
	// 	n := float64(i)
	// 	y := posY + (margin * (n + 1)) + (rowHeight * n)
	// 	ebitenutil.DrawRect(r, posX+margin, y, rowWidth, rowHeight, color.RGBA{50, 50, 50, 64})
	// }
}

// GetSize ...
func (s *BattleScene) GetSize() (int, int) {
	return s.size.Width, s.size.Height
}

// createCircleImage 半径rの円の画像イメージを作成します。color1は円の色、color2は円の向きを表す線の色です
func createRectImage(w, h int, color color.RGBA) *ebiten.Image {
	m := image.NewRGBA(image.Rect(0, 0, w, h))
	// 横ループ、半径*2＝直径まで
	for x := 0; x < w; x++ {
		// 縦ループ、半径*2＝直径まで
		for y := 0; y < h; y++ {
			m.Set(x, y, color)
		}
	}
	eimg, err := ebiten.NewImageFromImage(m, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	return eimg
}
