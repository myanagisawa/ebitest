package control

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/app/g"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/functions"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/models/event"
	"github.com/myanagisawa/ebitest/utils"
)

// HeaderBar ...
type HeaderBar struct {
	Base
	parent   interfaces.UIControl
	frameImg []*ebiten.Image
	closeBtn interfaces.UIControl
}

// NewHeaderBar ...
func NewHeaderBar(label string, parent interfaces.UIControl, draggable bool, hasCloseButton bool) interfaces.UIHeaderBar {
	f := make([]*ebiten.Image, 4)
	// フレームイメージを作成
	f[0] = ebiten.NewImageFromImage(utils.CreateRectImage(2, 2, &color.RGBA{220, 220, 220, 255}))
	f[1] = ebiten.NewImageFromImage(utils.CreateRectImage(2, 2, &color.RGBA{220, 220, 220, 255}))
	f[2] = ebiten.NewImageFromImage(utils.CreateRectImage(2, 2, &color.RGBA{220, 220, 220, 255}))
	f[3] = ebiten.NewImageFromImage(utils.CreateRectImage(2, 2, &color.RGBA{220, 220, 220, 255}))

	img := utils.CreateRectImage(1, 1, &color.RGBA{0, 0, 0, 192})

	b := &Base{
		label:        label,
		image:        ebiten.NewImageFromImage(img),
		position:     g.NewPoint(0, 0),
		scale:        g.NewScale(float64(parent.Size(enum.TypeScaled).W()), 20),
		visible:      true,
		eventHandler: event.NewEventHandler(),
	}
	o := &HeaderBar{
		Base:     *b,
		parent:   parent,
		frameImg: f,
	}
	if draggable {
		o.eventHandler.AddEventListener(enum.EventTypeDragging, functions.CommonEventCallback)
		o.eventHandler.AddEventListener(enum.EventTypeDragDrop, functions.CommonEventCallback)
	}

	if hasCloseButton {
		closeBtn := NewSimpleLabel("×", g.NewPoint(0, 0), 14, &color.RGBA{255, 255, 255, 255}, enum.FontStyleGenShinGothicBold)
		closeBtn.EventHandler().AddEventListener(enum.EventTypeFocus, functions.CommonEventCallback)
		closeBtn.EventHandler().AddEventListener(enum.EventTypeClick, func(o interfaces.EventOwner, pos *g.Point, params map[string]interface{}) {
			// 閉じる
			parent.SetVisible(false)
		})
		o.closeBtn = closeBtn
	}

	return o
}

// GetObjects ...
func (o *HeaderBar) GetObjects(x, y int) []interfaces.EbiObject {
	objs := []interfaces.EbiObject{}
	if o.closeBtn != nil && o.closeBtn.In(x, y) {
		objs = append(objs, o.closeBtn)
	}

	if o.In(x, y) {
		objs = append(objs, o)
	}
	return objs
}

// Draw draws the sprite.
func (o *HeaderBar) Draw(screen *ebiten.Image) {
	o.DrawWithOptions(screen, nil)
}

// DrawWithOptions draws the sprite.
func (o *HeaderBar) DrawWithOptions(screen *ebiten.Image, in *ebiten.DrawImageOptions) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	// 描画位置指定
	// log.Printf("in: %#v", in.GeoM.String())
	op.GeoM.Scale(o.Scale(enum.TypeGlobal).Get())
	op.GeoM.Translate(o.Position(enum.TypeGlobal).Get())
	// log.Printf("op: %#v / %#v", op, op.GeoM)

	if in != nil {
		op.GeoM.Concat(in.GeoM)
		op.ColorM.Concat(in.ColorM)
	}
	screen.DrawImage(o.image, op)

	// // 閉じるボタン
	if o.closeBtn != nil {
		op2 := &ebiten.DrawImageOptions{}

		op2.GeoM.Translate(o.Position(enum.TypeGlobal).Get())
		op2.GeoM.Translate(5, 0)

		if in != nil {
			op2.GeoM.Concat(in.GeoM)
		}
		o.closeBtn.DrawWithOptions(screen, op2)
	}

	// // 枠の描画
	o.drawFrameSet(screen)

	return op
}

func (o *HeaderBar) drawFrameSet(screen *ebiten.Image) {
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

}
