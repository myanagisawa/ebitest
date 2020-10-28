package scenes

import (
	"fmt"
	"image/color"
	"image/draw"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/myanagisawa/ebitest/utils"
	"golang.org/x/image/font"
)

// UIControl ...
type UIControl interface {
	Update() error
	Draw(screen *ebiten.Image)
	In(x, y int) bool
	SetLayer(l Layer)
}

// UIButton ...
type UIButton interface {
	UIControl
}

// UIControllerImpl ...
type UIControllerImpl struct {
	layer Layer
	image *ebiten.Image
	x     int
	y     int
}

// SetLayer ...
func (c *UIControllerImpl) SetLayer(l Layer) {
	c.layer = l
}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (c *UIControllerImpl) In(x, y int) bool {
	return c.image.At(x-c.x, y-c.y).(color.RGBA).A > 0
}

// Update ...
func (c *UIControllerImpl) Update() error {
	// log.Printf("UIControllerImpl: update")
	return nil
}

// Draw draws the sprite.
func (c *UIControllerImpl) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(c.image, op)
}

// UIButtonImpl ...
type UIButtonImpl struct {
	UIControllerImpl
	hover bool
}

// NewButton ...
func NewButton(label string, baseImg draw.Image, fontFace font.Face, labelColor color.Color, x, y int) UIButton {
	img := utils.SetTextToCenter(label, baseImg, fontFace, labelColor)
	eimg := ebiten.NewImageFromImage(*img)
	con := &UIControllerImpl{image: eimg, x: x, y: y}
	return &UIButtonImpl{UIControllerImpl: *con}
}

// Draw draws the sprite.
func (c *UIButtonImpl) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.x), float64(c.y))
	r, g, b, a := 1.0, 1.0, 1.0, 1.0
	if c.hover {
		r, g, b, a = 0.5, 0.5, 0.5, 1.0
	}
	op.ColorM.Scale(r, g, b, a)
	screen.DrawImage(c.image, op)
}

// UIScrollView ...
type UIScrollView interface {
	UIControl
}

// UIScrollViewImpl ...
type UIScrollViewImpl struct {
	UIControllerImpl
	contents        *ebiten.Image
	backgroundImage *ebiten.Image
	// imgparts        map[string]*ebiten.Image
	position      *Point
	contentScale  float64
	contentOffset *Point
	vScrollBar    *vScrollBar
	// scrollMin       float64
	// scrollMax       float64
	// scrollbarScale  float64
	// strokes         map[*Stroke]struct{}
}

type vScrollBar struct {
	imgparts       map[string]*ebiten.Image
	position       *Point
	draggingPos    *Point
	scrollMin      float64
	scrollMax      float64
	scrollbarScale float64
	strokes        map[*Stroke]struct{}
	hover          bool
}

// NewUIScrollView ...
func NewUIScrollView(x, y, w, h int, bgColor color.Color) UIScrollView {
	// 背景画像を作成（背景画像=表示領域）
	bgimg := createRectEbitenImage(w, h, bgColor)
	bw, _ := bgimg.Size()

	// スクロールパーツ作成
	scrollbaseimg := createRectEbitenImage(15, h, color.RGBA{255, 255, 255, 255})
	scrollbarimg := createRectEbitenImage(10, 10, color.RGBA{192, 192, 192, 255})
	scrollbarhilightimg := createRectEbitenImage(10, 10, color.RGBA{127, 127, 127, 255})

	img, _ := utils.GetImageByPath("resources/object_images/obj4.jpg")
	contents := ebiten.NewImageFromImage(img)

	cw, ch := contents.Size()
	wscale := float64(w-15) / float64(cw)
	hscale := float64(h*h-6) / (float64(ch) * wscale)

	s := &vScrollBar{
		imgparts: map[string]*ebiten.Image{
			"scrollBase":     scrollbaseimg,
			"scrollBar":      scrollbarimg,
			"scrollBarHover": scrollbarhilightimg,
		},
		position:       &Point{float64(bw) - 15, 0.0},
		draggingPos:    &Point{0, 0},
		scrollMin:      0,
		scrollMax:      float64(h) - (float64(ch) * wscale),
		scrollbarScale: hscale,
		strokes:        map[*Stroke]struct{}{},
	}

	c := &UIScrollViewImpl{
		backgroundImage: bgimg,
		contents:        contents,
		position:        &Point{float64(x), float64(y)},
		contentScale:    wscale,
		contentOffset:   &Point{0, 0},
		vScrollBar:      s,
		// imgparts: map[string]*ebiten.Image{
		// 	"scrollBase":     scrollbaseimg,
		// 	"scrollBar":      scrollbarimg,
		// 	"scrollBarHover": scrollbarhilightimg,
		// },
		// scrollMin:      0,
		// scrollMax:      float64(h) - (float64(ch) * wscale),
		// scrollbarScale: hscale,
		// strokes:        map[*Stroke]struct{}{},
	}
	return c
}

//
func createRectEbitenImage(w, h int, bgColor color.Color) *ebiten.Image {
	img := createRectImage(w, h, bgColor.(color.RGBA))
	eimg := ebiten.NewImageFromImage(img)
	return eimg
}

// Update ...
func (c *UIScrollViewImpl) Update() error {
	// ホイールイベント
	_, dy := ebiten.Wheel()
	// log.Printf("%0.1f < %0.1f && %0.1f > %0.1f", c.contentOffset.y, c.scrollMin, c.contentOffset.y, c.scrollMax)
	if dy < 0 {
		if c.contentOffset.y > c.vScrollBar.scrollMax {
			c.contentOffset.SetDelta(0, dy*5)
		}
	} else {
		if c.contentOffset.y < c.vScrollBar.scrollMin {
			c.contentOffset.SetDelta(0, dy*5)
		}
	}

	// ドラッグイベント
	for stroke := range c.vScrollBar.strokes {
		c.vScrollBar.updateStroke(stroke)
		if stroke.IsReleased() {
			delete(c.vScrollBar.strokes, stroke)
			log.Printf("drag end")
		}
	}

	// TODO: 共通化
	w, _ := c.backgroundImage.Size()
	cx, cy := ebiten.CursorPosition()
	lx, ly := c.layer.GetGlobalPosition()
	vx, vy := c.position.Get()
	ax := int(lx + vx + float64(w-15))
	ay := int(ly + vy)

	c.vScrollBar.hover = c.vScrollBar.imgparts["scrollBase"].At(cx-ax, cy-ay).(color.RGBA).A > 0
	if c.vScrollBar.hover {
		// log.Printf("hover")
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			stroke := NewStroke(&MouseStrokeSource{})
			// レイヤ内のドラッグ対象のオブジェクトを取得する仕組みが必要
			c.vScrollBar.strokes[stroke] = struct{}{}
			log.Printf("drag start")
		}
	}

	return nil
}

// Draw draws the sprite.
func (c *UIScrollViewImpl) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(c.position.Get())

	screen.DrawImage(c.backgroundImage, op)

	// // スクロールバー構成
	// w, h := c.backgroundImage.Size()
	// //スクロールエリア
	// op = &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(c.vScrollBar.position.Get())
	// c.backgroundImage.DrawImage(c.vScrollBar.imgparts["scrollBase"], op)

	// コンテンツ
	op = &ebiten.DrawImageOptions{}
	// log.Printf("wscale: %0.2f", wscale)
	op.GeoM.Scale(c.contentScale, c.contentScale)
	op.GeoM.Translate(c.contentOffset.Get())
	c.backgroundImage.DrawImage(c.contents, op)

	//スクロールバー
	_, ch := c.contents.Size()
	_, cy := c.contentOffset.Get()
	c.vScrollBar.Draw(c.backgroundImage, float64(ch), float64(cy), c.contentScale)
	// _, b := c.contentOffset.Get()
	// _, ch := c.contents.Size()
	// y := float64(h-6) * b / (float64(ch) * c.contentScale)
	// if math.Abs(y) < 3 {
	// 	y = -3
	// }
	// if math.Abs(y+c.vScrollBar.scrollbarScale) > float64(h-3) {
	// 	y = float64(h) - c.vScrollBar.scrollbarScale - 3
	// }
	// // log.Printf("表示高さ: %0.2f, Offset: %0.2f, コンテンツ高さ: %0.2f = バー移動量: %0.2f", float64(h), b, (float64(ch) * c.contentScale), y)
	// op = &ebiten.DrawImageOptions{}
	// op.GeoM.Scale(1.0, c.vScrollBar.scrollbarScale/10)
	// op.GeoM.Translate(float64(w-12), -y)

	// // TODO: 共通化
	// cx, cy := ebiten.CursorPosition()
	// lx, ly := c.layer.GetGlobalPosition()
	// vx, vy := c.position.Get()
	// ax := int(lx + vx + float64(w-15))
	// ay := int(ly + vy)
	// // log.Printf("cx=%d, cy=%d :: vx(%0.0f)+w(%d), vy(%0.0f)", cx, cy, lx+vx, w, ly+vy)

	// img := c.vScrollBar.imgparts["scrollBar"]
	// if c.vScrollBar.imgparts["scrollBase"].At(cx-ax, cy-ay).(color.RGBA).A > 0 {
	// 	img = c.vScrollBar.imgparts["scrollBarHover"]
	// }
	// c.backgroundImage.DrawImage(img, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("offset: %0.1f, %0.1f", c.contentOffset.x, c.contentOffset.y))
}

func (c *vScrollBar) updateStroke(stroke *Stroke) {
	stroke.Update()
	c.ScrollBy(stroke.PositionDiff())
}

// ScrollBy ...
func (c *vScrollBar) ScrollBy(x, y int) {
	// log.Printf("dragging: x=%d, y=%d", x, y)
	c.draggingPos.y = float64(y)
}

func (c *vScrollBar) Draw(r *ebiten.Image, contentHeight, contentOffsetY, contentScale float64) {
	rw, rh := r.Size()

	base := c.imgparts["scrollBase"]

	//スクロールエリア
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(c.position.Get())
	r.DrawImage(base, op)

	//スクロールバー
	op = &ebiten.DrawImageOptions{}
	translateY := (float64(rh) - 6) * contentOffsetY / (contentHeight * contentScale)
	// log.Printf("バー長さ: %0.2f, バー移動量: %0.2f", c.scrollbarScale, translateY)
	if math.Abs(translateY) < 3 {
		translateY = -3
	}
	if math.Abs(translateY+c.scrollbarScale) > (float64(rh) - 3) {
		translateY = float64(rh) - c.scrollbarScale - 3
	}
	// log.Printf("表示高さ: %0.2f, Offset: %0.2f, コンテンツ高さ: %0.2f = バー移動量: %0.2f", float64(rh), contentOffsetY, (contentHeight * contentScale), translateY)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(1.0, c.scrollbarScale/10)
	op.GeoM.Translate(float64(rw-12), -translateY)
	// op.GeoM.Translate(3.0, -translateY)

	bar := c.imgparts["scrollBar"]
	if c.hover {
		log.Printf("hover")
		bar = c.imgparts["scrollBarHover"]
	}
	r.DrawImage(bar, op)

}
