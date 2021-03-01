package control

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/app/g"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/functions"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/utils"
)

// closeBtn ...
type closeBtn struct {
	Base
	parent *HeaderBar
}

// newCloseBtn ...
func newCloseBtn(l interfaces.Layer, parent *HeaderBar) *closeBtn {
	b := NewSimpleLabel(l, "×", g.NewPoint(5.0, 0.0), 14, &color.RGBA{255, 255, 255, 255}, enum.FontStyleGenShinGothicBold).(*Base)
	o := &closeBtn{
		Base:   *b,
		parent: parent,
	}
	o.EventHandler().AddEventListener(enum.EventTypeFocus, functions.CommonEventCallback)
	o.EventHandler().AddEventListener(enum.EventTypeClick, func(eo interfaces.EventOwner, pos *g.Point, params map[string]interface{}) {
		// 閉じる
		parent.parent.SetVisible(false)
	})

	// o.SetPositionFunc(func(self interface{}, t enum.ValueTypeEnum) *g.Point {
	// 	o := self.(*closeBtn)
	// 	barPos := o.parent.Position(t)
	// 	return g.NewPoint(barPos.X()+o.position.X(), barPos.Y()+o.position.Y())
	// })
	return o
}

// In ...
func (o *closeBtn) In(x, y int) bool {
	if !o.visible {
		return false
	}

	return PointInBound(
		g.NewPoint(float64(x), float64(y)),
		g.NewBoundByPosSize(
			o.Position(enum.TypeGlobal),
			o.Size(enum.TypeScaled),
		),
		g.NewBoundByPosSize(
			o.Layer().Frame().Position(enum.TypeGlobal),
			o.Layer().Frame().Size(),
		),
	)
}

// HeaderBar ...
type HeaderBar struct {
	Base
	parent   interfaces.UIControl
	frameImg []*ebiten.Image
	closeBtn interfaces.UIControl
}

// NewHeaderBar ...
func NewHeaderBar(l interfaces.Layer, label string, parent interfaces.UIControl, draggable bool, hasCloseButton bool) interfaces.UIHeaderBar {
	f := make([]*ebiten.Image, 4)
	// フレームイメージを作成
	f[0] = ebiten.NewImageFromImage(utils.CreateRectImage(2, 2, &color.RGBA{220, 220, 220, 255}))
	f[1] = ebiten.NewImageFromImage(utils.CreateRectImage(2, 2, &color.RGBA{220, 220, 220, 255}))
	f[2] = ebiten.NewImageFromImage(utils.CreateRectImage(2, 2, &color.RGBA{220, 220, 220, 255}))
	f[3] = ebiten.NewImageFromImage(utils.CreateRectImage(2, 2, &color.RGBA{220, 220, 220, 255}))

	img := utils.CreateRectImage(1, 1, &color.RGBA{0, 0, 0, 192})

	b := NewControlBase(l, label, ebiten.NewImageFromImage(img), g.DefPoint(), g.NewScale(float64(parent.Size(enum.TypeScaled).W()), 20), true).(*Base)
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
		o.closeBtn = newCloseBtn(l, o)
	}

	// o.SetPositionFunc(func(self interface{}, t enum.ValueTypeEnum) *g.Point {
	// 	o := self.(*HeaderBar)
	// 	dlgPos := o.parent.Position(t)
	// 	return g.NewPoint(dlgPos.X()+o.position.X(), dlgPos.Y()+o.position.Y())
	// })

	l.AddUIControl(o)
	return o
}

// GetObjects ...
func (o *HeaderBar) GetObjects(x, y int) []interfaces.EbiObject {
	objs := []interfaces.EbiObject{}
	if o.closeBtn != nil && o.closeBtn.In(x, y) {
		log.Printf("くろーずぼたん: %d, %d", x, y)
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
	// // 枠の描画
	o.drawFrameSet(screen)
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
