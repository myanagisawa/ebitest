package control

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/example/t5/ebitest"
	"github.com/myanagisawa/ebitest/example/t5/interfaces"
	"github.com/myanagisawa/ebitest/example/t5/models"
	"github.com/myanagisawa/ebitest/utils"
	"golang.org/x/image/font"
)

// UIControlImpl ...
type UIControlImpl struct {
	label          string
	layer          interfaces.Layer
	bg             *models.EbiObject
	hasHoverAction bool
	hover          bool
}

// Label ...
func (c *UIControlImpl) Label() string {
	return fmt.Sprintf("%s.%s", c.layer.Label(), c.label)
}

// EbiObjects ...
func (c *UIControlImpl) EbiObjects() []*models.EbiObject {
	return []*models.EbiObject{c.bg}
}

// SetLayer ...
func (c *UIControlImpl) SetLayer(l interfaces.Layer) {
	c.layer = l
	c.bg.SetParent(l.EbiObjects()[0])
}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (c *UIControlImpl) In(x, y int) bool {
	// パーツ位置（左上座標）
	tx, ty := c.bg.GlobalPosition()

	// パーツサイズ(オリジナル)
	w, h := c.bg.EbitenImage().Size()

	// スケール
	sx, sy := c.bg.GlobalScale()

	// 見かけ上の右下座標を取得
	maxX := int(float64(w)*sx + tx)
	maxY := int(float64(h)*sy + ty)
	if maxX > ebitest.Width {
		maxX = ebitest.Width
	}
	if maxY > ebitest.Height {
		maxY = ebitest.Height
	}

	// 見かけ上の左上座標を取得
	minX, minY := int(tx), int(ty)
	if minX < 0 {
		minX = 0
	}
	if minY < 0 {
		minY = 0
	}
	// log.Printf("レイヤ座標: {(%d, %d), (%d, %d)}", minX, minY, maxX, maxY)
	return (x >= minX && x <= maxX) && (y > minY && y <= maxY)
	// return l.bg.At(x-l.x, y-l.y).(color.RGBA).A > 0
}

// HasHoverAction ...
func (c *UIControlImpl) HasHoverAction() bool {
	return c.hasHoverAction
}

// Update ...
func (c *UIControlImpl) Update() error {
	// log.Printf("UIControlImpl: update")
	if c.hasHoverAction {
		c.hover = c.In(ebiten.CursorPosition())
		if c.hover {
			log.Printf("hover: %s", c.label)
		}
	}
	return nil
}

// Draw draws the sprite.
func (c *UIControlImpl) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	w, h := c.bg.Size()
	// 描画位置指定
	op.GeoM.Reset()

	op.GeoM.Scale(c.bg.GlobalScale())

	// 対象画像の縦横半分だけマイナス位置に移動（原点に中心座標が来るように移動する）
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	// 中心を軸に回転
	op.GeoM.Rotate(c.bg.Theta())
	// ユニットの座標に移動
	op.GeoM.Translate(float64(w)/2, float64(h)/2)

	op.GeoM.Translate(c.bg.GlobalPosition())

	screen.DrawImage(c.bg.EbitenImage(), op)
}

// UIButtonImpl ...
type UIButtonImpl struct {
	UIControlImpl
}

// NewButton ...
func NewButton(label string, parent interfaces.Layer, baseImg draw.Image, fontFace font.Face, labelColor color.Color, x, y float64) interfaces.UIButton {
	img := utils.SetTextToCenter(label, baseImg, fontFace, labelColor)
	eimg := ebiten.NewImageFromImage(*img)

	con := &UIControlImpl{
		label:          label,
		bg:             models.NewEbiObject(label, eimg, parent.EbiObjects()[0], nil, ebitest.NewPoint(x, y), 0, true, true, false),
		hasHoverAction: true,
	}
	return &UIButtonImpl{UIControlImpl: *con}
}

// Draw draws the sprite.
func (c *UIButtonImpl) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	w, h := c.bg.Size()
	// 描画位置指定
	op.GeoM.Reset()

	op.GeoM.Scale(c.bg.GlobalScale())

	// 対象画像の縦横半分だけマイナス位置に移動（原点に中心座標が来るように移動する）
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	// 中心を軸に回転
	op.GeoM.Rotate(c.bg.Theta())
	// ユニットの座標に移動
	op.GeoM.Translate(float64(w)/2, float64(h)/2)

	op.GeoM.Translate(c.bg.GlobalPosition())
	r, g, b, a := 1.0, 1.0, 1.0, 1.0
	if c.hover {
		r, g, b, a = 0.5, 0.5, 0.5, 1.0
	}
	op.ColorM.Scale(r, g, b, a)
	screen.DrawImage(c.bg.EbitenImage(), op)
}

// UITextImpl ...
type UITextImpl struct {
	UIControlImpl
}

// NewText ...
func NewText(text string, parent interfaces.Layer, fontFace font.Face, c color.Color, x, y float64) interfaces.UIText {
	img := utils.CreateTextImage(text, fontFace, c)
	eimg := ebiten.NewImageFromImage(*img)

	label := fmt.Sprintf("text-%s", utils.RandomLC(8))
	con := &UIControlImpl{
		label: label,
		bg:    models.NewEbiObject(label, eimg, parent.EbiObjects()[0], nil, ebitest.NewPoint(x, y), 0, true, true, false),
	}
	return &UITextImpl{UIControlImpl: *con}
}

// Draw draws the sprite.
func (c *UITextImpl) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	w, h := c.bg.Size()
	// 描画位置指定
	op.GeoM.Reset()

	op.GeoM.Scale(c.bg.GlobalScale())

	// 対象画像の縦横半分だけマイナス位置に移動（原点に中心座標が来るように移動する）
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	// 中心を軸に回転
	op.GeoM.Rotate(c.bg.Theta())
	// ユニットの座標に移動
	op.GeoM.Translate(float64(w)/2, float64(h)/2)

	op.GeoM.Translate(c.bg.GlobalPosition())
	screen.DrawImage(c.bg.EbitenImage(), op)
}

// UIColumnImpl ...
type UIColumnImpl struct {
	UIControlImpl
	text *ebiten.Image
}

// NewColumn ...
func NewColumn(text string, parent interfaces.Layer, fontFace font.Face, labelColor color.Color, bgColor color.Color, x, y float64) interfaces.UIColumn {
	img := ebitest.CreateBorderedRectImage(500, 50, bgColor.(color.RGBA))
	eimg := ebiten.NewImageFromImage(img)

	label := fmt.Sprintf("col-%s", utils.RandomLC(8))
	con := &UIControlImpl{
		label:          label,
		bg:             models.NewEbiObject(label, eimg, parent.EbiObjects()[0], nil, ebitest.NewPoint(x, y), 0, true, true, false),
		hasHoverAction: true,
	}

	t := utils.CreateTextImage(text, fontFace, labelColor)
	timg := ebiten.NewImageFromImage(*t)

	return &UIColumnImpl{UIControlImpl: *con, text: timg}
}

// Draw draws the sprite.
func (c *UIColumnImpl) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	w, h := c.bg.Size()
	// 描画位置指定
	op.GeoM.Reset()

	op.GeoM.Scale(c.bg.GlobalScale())

	// 対象画像の縦横半分だけマイナス位置に移動（原点に中心座標が来るように移動する）
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	// 中心を軸に回転
	op.GeoM.Rotate(c.bg.Theta())
	// ユニットの座標に移動
	op.GeoM.Translate(float64(w)/2, float64(h)/2)

	op.GeoM.Translate(c.bg.GlobalPosition())

	r, g, b, a := 1.0, 1.0, 1.0, 1.0
	if c.hover {
		r, g, b, a = 0.5, 0.5, 0.5, 1.0
	}
	op.ColorM.Scale(r, g, b, a)

	screen.DrawImage(c.bg.EbitenImage(), op)

	c.DrawText(screen)
}

// DrawText draws the sprite.
func (c *UIColumnImpl) DrawText(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	w, h := c.text.Size()
	// 描画位置指定
	op.GeoM.Reset()

	op.GeoM.Scale(c.bg.GlobalScale())

	// 対象画像の縦横半分だけマイナス位置に移動（原点に中心座標が来るように移動する）
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	// 中心を軸に回転
	op.GeoM.Rotate(c.bg.Theta())
	// ユニットの座標に移動
	op.GeoM.Translate(float64(w)/2, float64(h)/2)

	bw, bh := c.bg.Size()
	_, by := c.bg.GlobalScale()
	a := float64(bh-h) * by / 2

	tx, ty := c.bg.GlobalPosition()

	op.GeoM.Translate(tx, ty+a)
	screen.DrawImage(c.text.SubImage(image.Rect(0, 0, bw, h)).(*ebiten.Image), op)
}
