package control

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/app/g"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/models/event"
	"github.com/myanagisawa/ebitest/utils"
)

// Dialog ...
type Dialog struct {
	Base
	items    []interfaces.UIControl
	frameImg []*ebiten.Image
}

// NewDialog ...
func NewDialog(label string, size *g.Size, visible bool) interfaces.UIDialog {
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

	b := &Base{
		label:        label,
		image:        ebiten.NewImageFromImage(img),
		position:     g.NewPoint(0, 0),
		scale:        g.NewScale(float64(size.W()), float64(size.H())),
		visible:      visible,
		eventHandler: event.NewEventHandler(),
	}
	o := &Dialog{
		Base:     *b,
		items:    []interfaces.UIControl{},
		frameImg: f,
	}
	return o
}

// Items ...
func (o *Dialog) Items() []interfaces.UIControl {
	return o.items
}

// SetItems ...
func (o *Dialog) SetItems(items []interfaces.UIControl) {
	o.items = items
}

// Draw draws the sprite.
func (o *Dialog) Draw(screen *ebiten.Image) {
	_ = o.DrawWithOptions(screen, nil)
}

// DrawWithOptions draws the sprite.
func (o *Dialog) DrawWithOptions(screen *ebiten.Image, op *ebiten.DrawImageOptions) *ebiten.DrawImageOptions {
	if !o.visible {
		// log.Printf("%sは不可視です。", o.label)
		return op
	}
	if op == nil {
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Reset()
	}

	// 描画位置指定
	op.GeoM.Scale(o.Scale(enum.TypeGlobal).Get())
	op.GeoM.Translate(o.Position(enum.TypeGlobal).Get())

	screen.DrawImage(o.image, op)

	// items 描画
	for i := range o.items {
		item := o.items[i]
		op.GeoM.Reset()
		op.GeoM.Translate(o.Position(enum.TypeGlobal).Get())
		op.ColorM.Reset()
		op.ColorM.Scale(0.7, 0.7, 0.7, 1.0)
		op = item.DrawWithOptions(screen, op)
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
