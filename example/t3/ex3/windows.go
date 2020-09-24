package ex3

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

type (
	// InfoScene ...
	InfoScene struct {
		size  *Size
		units []Unit
	}
)

var (
	posX, posY = 5.0, 35.0
	winWidth   = 400.0
	margin     = 3.0
	rowWidth   = winWidth - (margin * 2)
	rowHeight  = 30.0

	rownum    = 0
	winHeight = 0.0

	winImg image.Image

	textAdjust int
)

// NewInfoScene ...
func NewInfoScene(units []Unit) *InfoScene {
	rownum = len(units)
	winHeight = (rowHeight * float64(rownum)) + (margin * float64(rownum+1))

	winImg = createRectImage(int(winWidth), int(winHeight), color.RGBA{0, 0, 0, 64})

	textAdjust = int(rowHeight / 2)
	scene := &InfoScene{
		size:  &Size{Width: int(winWidth), Height: int(winHeight)},
		units: units,
	}
	return scene
}

// Update ...
func (s *InfoScene) Update() error {

	//log.Printf("BattleScene.Update")

	return nil
}

// Draw ...
func (s *InfoScene) Draw(r *ebiten.Image) {

	eimg, _ := ebiten.NewImageFromImage(winImg, ebiten.FilterDefault)
	for i := 0; i < rownum; i++ {
		n := float64(i)
		y := (margin * (n + 1)) + (rowHeight * n)

		adjust := textAdjust + (fface10White.uiFontMHeight / 2)

		colmargin := 0.0
		// 1. No
		colmargin = margin
		text.Draw(eimg, fmt.Sprintf("%d", i+1), fface10White.uiFont, int(colmargin), int(y)+adjust, fface10White.uiFontColor)

		// 2. unit画像
		colmargin += 30
		unit := s.units[i].(*UnitImpl)
		w, h := unit.GetEntity().image.Size()
		// 描画オプション: 中心基準に移動、中心座標で回転
		op := &ebiten.DrawImageOptions{}
		// 描画位置指定
		op.GeoM.Reset()
		// 対象画像の縦横半分だけマイナス位置に移動（原点に中心座標が来るように移動する）
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		// サイズを変更
		op.GeoM.Scale(0.5, 0.5)
		// ユニットの座標に移動
		op.GeoM.Translate(colmargin, y+(rowHeight/2))
		eimg.DrawImage(unit.GetEntity().image, op)

		// 3. unit名
		colmargin += 30
		fface := fface10White
		if unit.GetStatus() == -1 {
			fface = fface10Red
		}
		text.Draw(eimg, unit.Label, fface.uiFont, int(colmargin), int(y)+adjust, fface.uiFontColor)

		// 4. 現HP / 最大HP
		colmargin += 100
		text.Draw(eimg, fmt.Sprintf("%d / %d", unit.HP, unit.MaxHP), fface.uiFont, int(colmargin), int(y)+adjust, fface.uiFontColor)

	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(posX, posY)

	r.DrawImage(eimg, op)

	// ebitenutil.DrawRect(r, posX, posY, winWidth, winHeight, color.RGBA{0, 0, 0, 64})

	// for i := 0; i < 10; i++ {
	// 	n := float64(i)
	// 	y := posY + (margin * (n + 1)) + (rowHeight * n)
	// 	ebitenutil.DrawRect(r, posX+margin, y, rowWidth, rowHeight, color.RGBA{50, 50, 50, 64})
	// }
}

// GetSize ...
func (s *InfoScene) GetSize() (int, int) {
	return s.size.Width, s.size.Height
}

// createCircleImage 半径rの円の画像イメージを作成します。color1は円の色、color2は円の向きを表す線の色です
func createRectImage(w, h int, color color.RGBA) image.Image {
	m := image.NewRGBA(image.Rect(0, 0, w, h))
	// 横ループ、半径*2＝直径まで
	for x := 0; x < w; x++ {
		// 縦ループ、半径*2＝直径まで
		for y := 0; y < h; y++ {
			m.Set(x, y, color)
		}
	}
	return m
}
