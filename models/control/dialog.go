package control

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/app/g"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/utils"
)

// Dialog ...
type Dialog struct {
	Base
	items    []interfaces.UIControl
	frameImg []*ebiten.Image
}

// NewDialog ...
func NewDialog(l interfaces.Layer, label string, size *g.Size, visible bool) interfaces.UIDialog {
	f := make([]*ebiten.Image, 8)
	// フレームイメージを作成
	f[0] = ebiten.NewImageFromImage(utils.CreateRectImage(3, 3, &color.RGBA{220, 220, 220, 255}))
	f[1] = ebiten.NewImageFromImage(utils.CreateRectImage(3, 3, &color.RGBA{220, 220, 220, 255}))
	f[2] = ebiten.NewImageFromImage(utils.CreateRectImage(3, 3, &color.RGBA{220, 220, 220, 255}))
	f[3] = ebiten.NewImageFromImage(utils.CreateRectImage(3, 3, &color.RGBA{220, 220, 220, 255}))

	f[4] = ebiten.NewImageFromImage(utils.CreateRectImage(5, 5, &color.RGBA{192, 192, 192, 255}))
	f[5] = ebiten.NewImageFromImage(utils.CreateRectImage(5, 5, &color.RGBA{192, 192, 192, 255}))
	f[6] = ebiten.NewImageFromImage(utils.CreateRectImage(5, 5, &color.RGBA{192, 192, 192, 255}))
	f[7] = ebiten.NewImageFromImage(utils.CreateRectImage(5, 5, &color.RGBA{192, 192, 192, 255}))

	img := utils.CreateRectImage(1, 1, &color.RGBA{0, 0, 0, 127})

	b := NewControlBase(l, label, ebiten.NewImageFromImage(img), g.DefPoint(), g.NewScale(float64(size.W()), float64(size.H())), visible).(*Base)
	o := &Dialog{
		Base:     *b,
		items:    []interfaces.UIControl{},
		frameImg: f,
	}

	l.AddUIControl(o)
	return o
}

// GetObjects ...
func (o *Dialog) GetObjects(x, y int) []interfaces.EbiObject {
	objs := []interfaces.EbiObject{}
	for i := len(o.items) - 1; i >= 0; i-- {
		c := o.items[i]
		objs = append(objs, c.GetObjects(x, y)...)
	}

	if o.In(x, y) {
		objs = append(objs, o)
	}
	return objs
}

// Items ...
func (o *Dialog) Items() []interfaces.UIControl {
	return o.items
}

// AddItem ...
func (o *Dialog) AddItem(item interfaces.UIControl) {
	o.items = append(o.items, item)
}

// SetItems ...
func (o *Dialog) SetItems(items []interfaces.UIControl) {
	o.items = items
}

// Draw draws the sprite.
func (o *Dialog) Draw(screen *ebiten.Image) {
	o.DrawWithOptions(screen, nil)
}

// DrawWithOptions draws the sprite.
func (o *Dialog) DrawWithOptions(screen *ebiten.Image, in *ebiten.DrawImageOptions) *ebiten.DrawImageOptions {
	if !o.visible {
		// log.Printf("%sは不可視です。", o.label)
		return in
	}
	op := &ebiten.DrawImageOptions{}

	// 描画位置指定
	op.GeoM.Scale(o.Scale(enum.TypeGlobal).Get())
	op.GeoM.Translate(o.Position(enum.TypeGlobal).Get())

	if in != nil {
		op.GeoM.Concat(in.GeoM)
		op.ColorM.Concat(in.ColorM)
	}
	screen.DrawImage(o.image, op)

	// items 描画
	for i := range o.items {
		item := o.items[i]

		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Translate(o.Position(enum.TypeGlobal).Get())
		item.DrawWithOptions(screen, op2)

		log.Printf("item: %s", item.Label())
	}

	// 枠の描画
	o.drawFrameSet(screen)

	return op
}

func (o *Dialog) drawFrameSet(screen *ebiten.Image) {
	pos := o.Position(enum.TypeGlobal)
	bgSize := g.NewSize(o.image.Size())

	dlgScale := o.Scale(enum.TypeGlobal)
	dlgSize := g.NewSize(bgSize.W()*int(dlgScale.X()), bgSize.H()*int(dlgScale.Y()))

	// 枠線
	eimg := o.frameImg[0]
	imageSize := eimg.Bounds().Size()
	scaleX := float64(dlgSize.W()) / float64(imageSize.X)
	scaleY := 1.0
	drawFrame(screen, eimg, g.NewScale(scaleX, scaleY), pos, 0, 0)

	eimg = o.frameImg[1]
	imageSize = eimg.Bounds().Size()
	scaleX = 1.0
	scaleY = float64(dlgSize.H()) / float64(imageSize.Y)
	drawFrame(screen, eimg, g.NewScale(scaleX, scaleY), pos, float64(dlgSize.W())-float64(imageSize.X), 0)

	eimg = o.frameImg[2]
	imageSize = eimg.Bounds().Size()
	scaleX = 1.0
	scaleY = float64(dlgSize.H()) / float64(imageSize.Y)
	drawFrame(screen, eimg, g.NewScale(scaleX, scaleY), pos, 0, 0)

	eimg = o.frameImg[3]
	imageSize = eimg.Bounds().Size()
	scaleX = float64(dlgSize.W()) / float64(imageSize.X)
	scaleY = 1.0
	drawFrame(screen, eimg, g.NewScale(scaleX, scaleY), pos, 0, float64(dlgSize.H())-float64(imageSize.Y))

	// 角
	eimg = o.frameImg[4]
	drawFrame(screen, eimg, g.NewScale(1.0, 1.0), pos, -1, -1)

	eimg = o.frameImg[5]
	drawFrame(screen, eimg, g.NewScale(1.0, 1.0), pos, float64(dlgSize.W())-float64(eimg.Bounds().Size().X)+1, -1)

	eimg = o.frameImg[6]
	drawFrame(screen, eimg, g.NewScale(1.0, 1.0), pos, -1, float64(dlgSize.H())-float64(eimg.Bounds().Size().Y)+1)

	eimg = o.frameImg[7]
	drawFrame(screen, eimg, g.NewScale(1.0, 1.0), pos, float64(dlgSize.W())-float64(eimg.Bounds().Size().X)+1, float64(dlgSize.H())-float64(eimg.Bounds().Size().Y)+1)
}

func drawFrame(screen *ebiten.Image, img *ebiten.Image, scale *g.Scale, basePos *g.Point, x, y float64) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Reset()
	op.GeoM.Scale(scale.Get())
	op.GeoM.Translate(basePos.Get())
	op.GeoM.Translate(x, y)
	// log.Printf("1: %0.1f, %0.1f, op=%#v", scaleX, scaleY, op)
	screen.DrawImage(img, op)
}
