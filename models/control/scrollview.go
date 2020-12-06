package control

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/ebitest"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/models/char"
)

/*****************************************************************/

// UIScrollView ...
type UIScrollView struct {
	Base
	listView      *listView
	scrollbarBase *scrollbarBase
	fontSet       *char.Resource
}

// SetDataSource ...
func (o *UIScrollView) SetDataSource(colNames []interface{}, data [][]interface{}) {
	// まずは幅を測ってカラムサイズを決めてしまうのが良いか。

	o.listView.SetRows(data)
}

// Update ...
func (o *UIScrollView) Update() error {
	o.listView.Update()
	o.scrollbarBase.Update()

	return nil
}

// Draw ...
func (o *UIScrollView) Draw(screen *ebiten.Image) {
	o.Base.Draw(screen)
	o.listView.Draw(screen)
	o.scrollbarBase.Draw(screen)
}

// listViewSize スクロールビューの中のリスト領域のサイズを返す
func (o *UIScrollView) listViewSize() (int, int) {
	iw, ih := o.image.Size()

	//リスト表示領域の調整をここで実施
	w, h := iw, ih-20
	return w, h
}

// NewUIScrollView ...
func NewUIScrollView(label string, pos *ebitest.Point, size *ebitest.Size) interfaces.UIScrollView {

	// スクロールビュー全体のベース画像
	img := ebitest.CreateRectImage(size.W(), size.H(), &color.RGBA{255, 0, 0, 127})
	eimg := ebiten.NewImageFromImage(img)

	// eimg := ebiten.NewImage(size.Get())

	cb := Base{
		label:          label,
		image:          eimg,
		position:       pos,
		scale:          ebitest.NewScale(1.0, 1.0),
		hasHoverAction: false,
	}

	o := &UIScrollView{
		Base:    cb,
		fontSet: char.Res.Get(12, enum.FontStyleGenShinGothicNormal),
	}

	// スクロール部分の初期化
	list := newListView(fmt.Sprintf("%s.list", label), o, ebitest.NewPoint(0, 20))
	o.listView = list

	// スクロールバー部分の初期化
	barBase := newScrollbarBase(fmt.Sprintf("%s.scrollbar.base", label), o, ebitest.NewPoint(float64(size.W()-15), 20))
	o.scrollbarBase = barBase

	return o
}

/*****************************************************************/

// scrollViewParts パーツの基底
type scrollViewParts struct {
	Base
	parent *UIScrollView
}

// Position ...
func (o *scrollViewParts) Position(t enum.ValueTypeEnum) *ebitest.Point {
	if t == enum.TypeLocal {
		return o.position
	}
	gx, gy := 0.0, 0.0
	if o.parent != nil {
		gx, gy = o.parent.Position(enum.TypeGlobal).Get()
	}
	sx, sy := o.Scale(enum.TypeGlobal).Get()
	gx += o.position.X() * sx
	gy += o.position.Y() * sy
	return ebitest.NewPoint(gx, gy)
}

// Draw ...
func (o *scrollViewParts) Draw(screen *ebiten.Image) {
	var op *ebiten.DrawImageOptions

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(o.Position(enum.TypeGlobal).Get())

	screen.DrawImage(o.image, op)
}

/*****************************************************************/

// listView ...
type listView struct {
	scrollViewParts
	scrollingPos *ebitest.Point
	rows         []*listRow
}

// SetRows ...
func (o *listView) SetRows(data [][]interface{}) {

	rows := make([]*listRow, len(data))
	for i := range data {
		rows[i] = newListRow(fmt.Sprintf("listRow[%d]", i), o.parent, i, data[i])
	}
	o.rows = rows
}

// Update ...
func (o *listView) Update() error {
	// ホイールイベント
	_, dy := ebiten.Wheel()
	o.scrollingPos.SetDelta(0, dy*2)
	if o.scrollingPos.Y() < 0 {
		// 上に余白ができる
		o.scrollingPos.Set(0, 0)
	} else {
		_, ih := o.image.Size()
		_, ph := o.parent.listViewSize()
		if int(o.scrollingPos.Y())+ph > ih {
			// 下に余白ができる
			o.scrollingPos.Set(0, float64(ih-ph))
		}
	}

	return nil
}

// Draw ...
func (o *listView) Draw(screen *ebiten.Image) {
	var op *ebiten.DrawImageOptions

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(o.Position(enum.TypeGlobal).Get())

	_, sy := o.scrollingPos.GetInt()
	w, h := o.parent.listViewSize()
	fr := image.Rect(0, sy, w, h+sy)
	// log.Printf("%s: pos: %0.1f, %0.1f, fr: %d, %d, %d, %d", o.label, o.Position(enum.TypeGlobal).X(), o.Position(enum.TypeGlobal).Y(), fr.Min.X, fr.Min.Y, fr.Max.X, fr.Max.Y)
	screen.DrawImage(o.image.SubImage(fr).(*ebiten.Image), op)

	// rows描画
	x, y := o.Position(enum.TypeGlobal).GetInt()
	th := 0
	for i := range o.rows {
		row := o.rows[i]
		tw, nw, nh := 0, 0, 0
		str := fmt.Sprintf("ty=%d, y=%d, h=%d", int(y+th-sy), y, h)
		imgs := o.parent.fontSet.GetByString(str)
		list := append(row.texts, imgs...)
		for j := range list {
			nw, nh = list[j].Size()

			tx := float64(x + tw)
			ty := float64(y + th - sy)
			// text shadow
			ops := &ebiten.DrawImageOptions{}
			ops.GeoM.Translate(tx+1, ty+1)
			ops.ColorM.Scale(0, 0, 0, 0.5)
			screen.DrawImage(list[j], ops)

			op = &ebiten.DrawImageOptions{}
			op.GeoM.Translate(tx, ty)
			if int(ty)+nh < y || int(ty) > y+h {
				// リスト表示領域外
				op.ColorM.Scale(0.3, 0.3, 0.3, 0.5)
			}
			screen.DrawImage(list[j], op)

			tw += nw
		}
		th += nh
	}

}

func newListView(label string, parent *UIScrollView, pos *ebitest.Point) *listView {
	img := ebitest.Images["world"]
	eimg := ebiten.NewImageFromImage(img)

	// positionは親positionからのdeltaを指定する
	cb := Base{
		label:          label,
		image:          eimg,
		position:       pos,
		scale:          ebitest.NewScale(1.0, 1.0),
		hasHoverAction: false,
	}

	o := &listView{
		scrollViewParts: scrollViewParts{
			Base:   cb,
			parent: parent,
		},
		scrollingPos: ebitest.NewPoint(0, 0),
	}
	return o
}

/*****************************************************************/

// listRow ...
type listRow struct {
	scrollViewParts
	index  int
	texts  []*ebiten.Image
	source []interface{}
}

func newListRow(label string, parent *UIScrollView, idx int, row []interface{}) *listRow {
	cb := Base{
		label:          label,
		scale:          ebitest.NewScale(1.0, 1.0),
		hasHoverAction: true,
	}

	o := &listRow{
		scrollViewParts: scrollViewParts{
			Base:   cb,
			parent: parent,
		},
		index:  idx,
		source: row,
	}

	// テキスト画像生成
	for i := range row {
		str := fmt.Sprintf("%v / ", row[i])
		col := parent.fontSet.GetByString(str)
		o.texts = append(o.texts, col...)
	}

	return o
}

/*****************************************************************/

// scrollbarBase ...
type scrollbarBase struct {
	scrollViewParts
}

func newScrollbarBase(label string, parent *UIScrollView, pos *ebitest.Point) *scrollbarBase {
	_, ph := parent.listViewSize()
	scrollbaseimg := ebitest.CreateRectImage(15, ph, &color.RGBA{255, 255, 255, 64})
	eimg := ebiten.NewImageFromImage(scrollbaseimg)

	// positionは親positionからのdeltaを指定する
	cb := Base{
		label:          label,
		image:          eimg,
		position:       pos,
		scale:          ebitest.NewScale(1.0, 1.0),
		hasHoverAction: false,
	}

	o := &scrollbarBase{
		scrollViewParts: scrollViewParts{
			Base:   cb,
			parent: parent,
		},
	}
	return o
}
