package scenes

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
)

// Layer ...
type Layer interface {
	Label() string
	In(x, y int) bool
	GetGlobalPosition() (float64, float64)
	AddUIControl(c UIControl)
	UIControlAt(x, y int) UIControl
	ScaleTo(scale float64)
	LocalPosition(x, y int) (int, int)
	IsModal() bool
	AddEventListener(c UIControl, name string, callback func(UIControl, *EventSource))
	FiringEvent(name string)
	Update(screen *ebiten.Image) error
	Draw(screen *ebiten.Image)
}

// LayerBase ...
type LayerBase struct {
	label        string
	bg           *ebiten.Image
	x            int
	y            int
	scale        float64
	translateX   float64
	translateY   float64
	parent       Scene
	isModal      bool
	controls     []UIControl
	eventHandler *EventHandler
}

// Label ...
func (l *LayerBase) Label() string {
	return l.label
}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (l *LayerBase) In(x, y int) bool {
	// レイヤ位置（左上座標）
	tx, ty := int(l.translateX), int(l.translateY)

	// レイヤサイズ(オリジナル)
	w, h := l.bg.Size()

	// 見かけ上の右下座標を取得
	maxX := int(float64(w)*l.scale) + tx
	maxY := int(float64(h)*l.scale) + ty
	if maxX > width {
		maxX = width
	}
	if maxY > height {
		maxY = height
	}

	// 見かけ上の左上座標を取得
	minX, minY := tx, ty
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

// GetGlobalPosition ...
func (l *LayerBase) GetGlobalPosition() (float64, float64) {
	return l.translateX, l.translateY
}

// AddUIControl レイヤに部品を追加します
func (l *LayerBase) AddUIControl(c UIControl) {
	c.SetLayer(l)
	l.controls = append(l.controls, c)
}

// UIControlAt (x, y)座標に存在する部品を返します
func (l *LayerBase) UIControlAt(x, y int) UIControl {
	for i := len(l.controls) - 1; i >= 0; i-- {
		c := l.controls[i]
		if c.In(x, y) {
			return c
		}
	}
	return nil
}

// ScaleTo ...
func (l *LayerBase) ScaleTo(scale float64) {
	l.scale = scale

	w, h := l.bg.Size()
	w = int(float64(w) * l.scale)
	h = int(float64(h) * l.scale)

	if l.x < width-w {
		l.x = width - w
	}
	if l.y < height-h {
		l.y = height - h
	}
	l.translateX = float64(l.x)
	l.translateY = float64(l.y)
	// log.Printf("MoveBy: s.x=%0.2f, s.y=%0.2f", s.translateX, s.translateY)
}

// LocalPosition スクリーン上の座標をシーンオブジェクト上の座標に変換します
func (l *LayerBase) LocalPosition(x, y int) (int, int) {
	sx := float64(l.x) * l.scale * -1
	sy := float64(l.y) * l.scale * -1
	localX := int((float64(x) + sx) / l.scale)
	localY := int((float64(y) + sy) / l.scale)

	// log.Printf("scale: %0.2f [x: %d, s.x: %d = %d] [y: %d, s.y: %d = %d]", s.scale, x, s.x, localX, y, s.y, localY)
	return localX, localY
}

// IsModal ...
func (l *LayerBase) IsModal() bool {
	return l.isModal
}

// AddEventListener ...
func (l *LayerBase) AddEventListener(c UIControl, name string, callback func(UIControl, *EventSource)) {
	ev := &Event{c, callback}
	l.eventHandler.Set(name, ev)
}

// FiringEvent ...
func (l *LayerBase) FiringEvent(name string) {
	x, y := l.LocalPosition(ebiten.CursorPosition())
	l.eventHandler.Firing(l.parent, name, x, y)
}

// Update ...
func (l *LayerBase) Update(screen *ebiten.Image) error {
	if l.parent.ActiveLayer() == nil {
		return nil
	}
	// log.Printf("LayerBase.Update(): %s, l: %s", l.parent.ActiveLayer().Label(), l.label)
	if l.parent.ActiveLayer().Label() == l.label {
		// log.Printf("LayerBase.Update()")
		for _, c := range l.controls {

			_ = c.Update(screen)
		}
	}

	return nil
}

// Draw ...
func (l *LayerBase) Draw(screen *ebiten.Image) {
	// log.Printf("LayerBase.Draw")
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(l.translateX, l.translateY)

	screen.DrawImage(l.bg, op)

	for _, c := range l.controls {
		c.Draw(l.bg)
	}
}

// NewLayerBase ...
func NewLayerBase(parent Scene) *LayerBase {
	img := createRectImage(width, height, color.RGBA{32, 32, 32, 255})
	eimg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	l := &LayerBase{
		label:    "layer base",
		bg:       eimg,
		x:        0,
		y:        0,
		scale:    1.0,
		parent:   parent,
		isModal:  false,
		controls: []UIControl{},
		eventHandler: &EventHandler{
			events: map[string]map[*Event]struct{}{},
		},
	}

	l.translateX = float64(l.x)
	l.translateY = float64(l.y)
	return l
}

// TestWindow ...
type TestWindow struct {
	LayerBase
}

// NewTestWindow ...
func NewTestWindow(parent Scene) *TestWindow {
	img := createRectImage(240, 300, color.RGBA{0, 0, 0, 128})
	eimg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	l := &TestWindow{
		LayerBase: LayerBase{
			label:    "test window",
			bg:       eimg,
			x:        25,
			y:        100,
			scale:    1.0,
			parent:   parent,
			isModal:  false,
			controls: []UIControl{},
			eventHandler: &EventHandler{
				events: map[string]map[*Event]struct{}{},
			},
		},
	}
	l.translateX = float64(l.x)
	l.translateY = float64(l.y)

	c := NewButton("ユニット一覧", images["btnBase"], fonts["btnFont"], color.Black, 20, 20)

	l.AddUIControl(c)
	l.AddEventListener(c, "click", func(target UIControl, source *EventSource) {
		log.Printf("Open Sub clicked: x=%d, y=%d", source.x, source.y)
		log.Printf("target: %#v", target)

		source.scene.SetLayer(NewUnitListWindow(source.scene))
	})
	log.Printf("controls: %#v", l.controls)

	return l
}

// Draw ...
func (l *TestWindow) Draw(screen *ebiten.Image) {
	// log.Printf("TestWindow.Draw")

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(l.translateX, l.translateY)

	screen.DrawImage(l.bg, op)

	for _, c := range l.controls {
		c.Draw(l.bg)
	}
}

// TestSubWindow ...
type TestSubWindow struct {
	LayerBase
	coverimg *ebiten.Image
}

// NewTestSubWindow ...
func NewTestSubWindow(parent Scene) *TestSubWindow {
	img := createRectImage(600, 400, color.RGBA{0, 0, 0, 128})
	eimg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	l := &TestSubWindow{
		LayerBase: LayerBase{
			label:    "test sub window",
			bg:       eimg,
			x:        100,
			y:        200,
			scale:    1.0,
			parent:   parent,
			isModal:  true,
			controls: []UIControl{},
			eventHandler: &EventHandler{
				events: map[string]map[*Event]struct{}{},
			},
		},
	}
	l.translateX = float64(l.x)
	l.translateY = float64(l.y)

	subimg := createRectImage(width, height, color.RGBA{32, 32, 32, 192})
	subeimg, _ := ebiten.NewImageFromImage(subimg, ebiten.FilterDefault)

	l.coverimg = subeimg

	c := NewButton("閉じる", images["btnBase"], fonts["btnFont"], color.Black, 200, 200)
	l.AddUIControl(c)
	l.AddEventListener(c, "click", func(target UIControl, source *EventSource) {
		log.Printf("閉じる clicked: x=%d, y=%d", source.x, source.y)
		log.Printf("target: %#v", target)

		source.scene.DeleteLayer(source.scene.ActiveLayer())
	})

	return l
}

// Draw ...
func (l *TestSubWindow) Draw(screen *ebiten.Image) {
	// log.Printf("TestWindow.Draw")

	// modal背景を描画
	cop := &ebiten.DrawImageOptions{}
	cop.GeoM.Translate(0, 0)

	screen.DrawImage(l.coverimg, cop)

	// layer描画
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(l.translateX, l.translateY)

	screen.DrawImage(l.bg, op)

	for _, c := range l.controls {
		c.Draw(l.bg)
	}
}

// UnitListWindow ...
type UnitListWindow struct {
	LayerBase
	coverimg *ebiten.Image
}

// NewUnitListWindow ...
func NewUnitListWindow(parent Scene) *UnitListWindow {
	img := createRectImage(1000, 600, color.RGBA{0, 0, 0, 128})
	eimg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	l := &UnitListWindow{
		LayerBase: LayerBase{
			label:    "Unit List",
			bg:       eimg,
			x:        100,
			y:        100,
			scale:    1.0,
			parent:   parent,
			isModal:  true,
			controls: []UIControl{},
			eventHandler: &EventHandler{
				events: map[string]map[*Event]struct{}{},
			},
		},
	}
	l.translateX = float64(l.x)
	l.translateY = float64(l.y)

	// ScrollView
	list := NewUIScrollView(50, 50, 900, 450, color.RGBA{32, 32, 32, 192})
	l.AddUIControl(list)

	subimg := createRectImage(width, height, color.RGBA{32, 32, 32, 192})
	subeimg, _ := ebiten.NewImageFromImage(subimg, ebiten.FilterDefault)
	l.coverimg = subeimg

	// listBase, _ := ebiten.NewImageFromImage(images["listBase"], ebiten.FilterDefault)
	// l.listBase = listBase

	// listScroller, _ := ebiten.NewImageFromImage(images["listScroller"], ebiten.FilterDefault)
	// l.listScroller = listScroller

	// c := NewButton("閉じる", images["btnBase"], fonts["btnFont"], color.Black, 780, 540)
	// l.AddUIControl(c)
	// l.AddEventListener(c, "click", func(target UIControl, source *EventSource) {
	// 	log.Printf("閉じる clicked: x=%d, y=%d", source.x, source.y)
	// 	log.Printf("target: %#v", target)

	// 	source.scene.DeleteLayer(source.scene.ActiveLayer())
	// })

	return l
}

// Draw ...
func (l *UnitListWindow) Draw(screen *ebiten.Image) {
	// log.Printf("TestWindow.Draw")

	// modal背景を描画
	cop := &ebiten.DrawImageOptions{}
	cop.GeoM.Translate(0, 0)

	screen.DrawImage(l.coverimg, cop)

	// layer描画
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(l.translateX, l.translateY)

	screen.DrawImage(l.bg, op)

	// // listbase描画
	// l.DrawList(l.listBase)

	// op = &ebiten.DrawImageOptions{}
	// op.GeoM.Scale(90.0, 45.0)
	// op.GeoM.Translate(50, 50)

	// l.bg.DrawImage(l.listBase, op)

	for _, c := range l.controls {
		c.Draw(l.bg)
	}
}

// // DrawList ...
// func (l *UnitListWindow) DrawList(listBase *ebiten.Image) {
// 	// log.Printf("TestWindow.Draw")

// 	// listBase描画
// 	op := &ebiten.DrawImageOptions{}
// 	op.GeoM.Scale(50.0, 120.0)
// 	op.GeoM.Translate(7.0, 0.0)

// 	listBase.DrawImage(l.listScroller, op)
// }
