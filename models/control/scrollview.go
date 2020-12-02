package control

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/myanagisawa/ebitest/ebitest"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/example/t5/models"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/models/input"
	"github.com/myanagisawa/ebitest/utils"
)

const (
	margin  = 2
	padleft = 2

	sbWidth  = 15
	sbHeight = 10
)

var (
	scrollbaseimg, scrollbarimg, scrollbarhilightimg image.Image
)

type listRowView struct {
	Base
	source   []interface{}
	texts    []*ebiten.Image
	cols     []*ebiten.Image
	colWidth []float64
}

func init() {
	// スクロールパーツ作成
	scrollbaseimg = ebitest.CreateRectImage(sbWidth, sbHeight, color.RGBA{255, 255, 255, 255})
	scrollbarimg = ebitest.CreateRectImage(sbWidth-5, sbHeight, color.RGBA{192, 192, 192, 255})
	scrollbarhilightimg = ebitest.CreateRectImage(sbWidth-5, sbHeight, color.RGBA{127, 127, 127, 255})
}

// newListRowView スクロールリストの行オブジェクトを作成
func newListRowView(row []interface{}, w, h int, pos *ebitest.Point) *listRowView {

	img := ebitest.CreateRectImage(w, h, color.RGBA{0, 0, 0, 255})
	eimg := ebiten.NewImageFromImage(img)
	// eimg := ebiten.NewImage(w, h)

	label := fmt.Sprintf("row-%s", utils.RandomLC(8))
	con := &Base{
		label:          label,
		image:          eimg,
		position:       pos,
		hasHoverAction: true,
	}

	texts := []*ebiten.Image{}
	var timg *ebiten.Image
	for i := range row {
		col := row[i]
		switch val := col.(type) {
		case string:
			t := utils.CreateTextImage(val, ebitest.Fonts["btnFont"], color.RGBA{32, 32, 32, 255})
			timg = ebiten.NewImageFromImage(*t)
			texts = append(texts, timg)
			// log.Printf("string: %s", val)
		case int:
			t := utils.CreateTextImage(fmt.Sprintf("%d", val), ebitest.Fonts["btnFont"], color.RGBA{32, 32, 32, 255})
			timg = ebiten.NewImageFromImage(*t)
			texts = append(texts, timg)
			// log.Printf("int: %d", val)
		}
	}

	return &listRowView{
		Base:   *con,
		source: row,
		texts:  texts,
		cols:   []*ebiten.Image{},
	}
}

func (o *listRowView) GetRowHeight() int {
	_, h := o.image.Size()
	return h
}

// DrawListRow ...
func (o *listRowView) DrawListRow(hover bool) *ebiten.Image {
	var op *ebiten.DrawImageOptions

	bgSize := ebitest.NewSize(o.image.Size())
	base := ebiten.NewImage(bgSize.Get())

	x := 0.0
	h := 0
	for i, row := range o.texts {
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, margin)
		if hover {
			op.ColorM.Scale(0.5, 0.5, 0.5, 1.0)
		}
		base.DrawImage(o.cols[i], op)

		_, h = row.Size()

		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x+padleft, float64(bgSize.H()-h)/2)
		base.DrawImage(row, op)

		x += o.colWidth[i] + margin
	}

	return base
}

type listView struct {
	Base
	names *listRowView
	rows  []listRowView
}

// newListView ...
func newListView(colNames []interface{}, data [][]interface{}, roww, rowh int, pos *ebitest.Point) *listView {
	listh := 0
	rows := []listRowView{}
	for i := range data {
		rowview := newListRowView(data[i], roww, rowh, ebitest.NewPoint(0, float64(listh)))
		rows = append(rows, *rowview)
		listh += rowh
	}
	// margin分を高さに追加
	listh += margin * (len(data) - 1)

	// 見出し行作成
	names := newListRowView(colNames, roww, rowh, ebitest.NewPoint(0, 0))

	// 全データから列のサイズ比を取得
	if len(rows) > 0 {
		// 列ごとの最大幅を取得
		maxWidths := make([]int, len(rows[0].texts))
		for i := range rows {
			row := rows[i]
			for j := range row.texts {
				text := row.texts[j]
				w, _ := text.Size()
				if w > maxWidths[j] {
					maxWidths[j] = w
				}
			}
		}
		// 列名行も対象
		for i := range names.texts {
			text := names.texts[i]
			w, _ := text.Size()
			if w > maxWidths[i] {
				maxWidths[i] = w
			}
		}

		// 最大幅での各列のサイズ比を計算
		totalWidth := 0.0
		ratio := make([]float64, len(maxWidths))
		for i := range maxWidths {
			totalWidth += float64(maxWidths[i])
		}
		for i := range maxWidths {
			ratio[i] = float64(maxWidths[i]) / totalWidth
		}

		// カラムサイズリストを取得
		colWidth := make([]float64, len(ratio))
		for i := range ratio {
			colWidth[i] = float64(roww)*ratio[i] - margin
		}

		// カラム背景画像を作成
		colBg := make([]*ebiten.Image, len(colWidth))
		for i := range colWidth {
			img := ebitest.CreateRectImage(int(colWidth[i]), rowh-margin, color.RGBA{64, 64, 64, 64})
			eimg := ebiten.NewImageFromImage(img)
			colBg[i] = eimg
		}

		// rowViewにカラムサイズリストをセット
		for i := range rows {
			row := rows[i]
			row.colWidth = colWidth
			row.cols = colBg
			rows[i] = row
		}
		// 列名にもセット
		names.colWidth = colWidth
		names.cols = colBg

	}

	// ベース画像を作成
	eimg := ebiten.NewImage(roww, listh)

	label := fmt.Sprintf("list-%s", utils.RandomLC(8))
	con := &Base{
		label:          label,
		image:          eimg,
		position:       pos,
		hasHoverAction: true,
	}

	l := &listView{
		Base:  *con,
		names: names,
		rows:  rows,
	}
	return l
}

// DrawList ...
func (o *listView) DrawList(drawRect image.Rectangle) *ebiten.Image {
	var op *ebiten.DrawImageOptions
	eimg := o.image
	bw, bh := eimg.Size()

	base := ebiten.NewImage(bw, bh)

	// カーソルIN判定
	isHover := false
	cx, cy := ebiten.CursorPosition()
	bgPos := o.Position(enum.TypeGlobal)
	viewSize := ebitest.NewSize(drawRect.Dx(), drawRect.Dy())
	_, dy := 0, 0
	if int(bgPos.X()) <= cx && int(bgPos.X())+viewSize.W() >= cx {
		// 横位置がスクロールビュー範囲内
		if int(bgPos.Y()) <= cy && int(bgPos.Y())+viewSize.H() >= cy {
			// 縦位置がスクロールビュー範囲内
			// dx = cx - int(bgPos.X())
			dy = cy - int(bgPos.Y())
			isHover = true
		}
	}

	y := 0.0
	top, bottom, min, max := float64(drawRect.Min.Y), float64(drawRect.Max.Y), 0.0, 0.0
	for _, row := range o.rows {
		min = y
		max = y + float64(row.GetRowHeight())

		// 対象行の下端がtop以下あるいは対象行の上端がbottom以上、以外が描画対象
		if !(max <= top || min >= bottom) {
			op = &ebiten.DrawImageOptions{}
			op.GeoM.Translate(0, y)
			hov := false
			if isHover {
				if int(min-top) <= dy && int(max-top) >= dy {
					// log.Printf("カーソルは%d行目の範囲内: x=%d, y=%d", i, dx, dy)
					hov = true
				}
			}
			base.DrawImage(row.DrawListRow(hov), op)
		}

		y += float64(row.GetRowHeight()) + margin
	}

	return base
}

func (o *listView) GetListHeight() int {
	listHeight := 0
	for i := range o.rows {
		row := o.rows[i]
		listHeight += row.GetRowHeight()
	}
	return listHeight
}

// scrollBar ...
type scrollBar struct {
	base        *ebiten.Image
	bar         *ebiten.Image
	barHover    *ebiten.Image
	parent      *UIScrollViewImpl
	position    *ebitest.Point
	draggingPos *ebitest.Point
	scrollMin   float64
	scrollMax   float64
	stroke      *input.Stroke
	// hover       bool
	// scrollbarScale float64
}

func (o *scrollBar) updateStroke(stroke *input.Stroke) {
	stroke.Update()
	o.ScrollBy(stroke.PositionDiff())
}

// ScrollBy ...
func (o *scrollBar) ScrollBy(x, y float64) {
	// log.Printf("dragging: x=%d, y=%d", x, y)
	bgSize := ebitest.NewSize(o.base.Size())
	barSize := ebitest.NewSize(o.bar.Size())
	// log.Printf("ScrollBy: bgSize: %#v", bgSize)
	ratio := float64(barSize.H()) / float64(bgSize.H())
	o.draggingPos.Set(o.draggingPos.X(), y/ratio)
}

func (o *scrollBar) Draw(r *ebiten.Image, contentHeight, contentOffsetY, contentScale float64) {
	// スクロールビューサイズ
	baseSize := ebitest.NewSize(o.base.Size())
	// スクロールバーベース位置
	basePos := o.base.Bounds().Min
	// 親位置
	parentPos := c.base.Parent().GlobalPosition()
	// スクロールバーサイズ
	barSize := c.bar.Size()

	//スクロールエリア
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(c.base.Scale().X(), c.base.Scale().Y())
	op.GeoM.Translate(parentPos.X()+basePos.X(), parentPos.Y()+basePos.Y())
	r.DrawImage(c.base.EbitenImage(), op)

	//スクロールバー
	op = &ebiten.DrawImageOptions{}
	translateY := (float64(baseSize.H()) - 6) * contentOffsetY / (contentHeight * contentScale)
	// log.Printf("バー長さ: %0.2f, バー移動量: %0.2f", c.scrollbarScale, translateY)
	if math.Abs(translateY) < 3 {
		translateY = 3
	}

	if math.Abs(translateY+float64(barSize.H())) > (float64(baseSize.H()) - 3) {
		translateY = float64(baseSize.H() - barSize.H() - 3)
	}
	// log.Printf("表示高さ: %0.2f, Offset: %0.2f, コンテンツ高さ: %0.2f = バー移動量: %0.2f", float64(rh), contentOffsetY, (contentHeight * contentScale), translateY)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(c.bar.Scale().X(), c.bar.Scale().Y())
	op.GeoM.Translate(parentPos.X()+basePos.X()+3, parentPos.Y()+basePos.Y()+translateY)
	// op.GeoM.Translate(3.0, -translateY)
	// log.Printf("op: %#v", op)
	bar := c.bar.EbitenImage()
	if c.base.In(ebiten.CursorPosition()) {
		// log.Printf("hover")
		bar = c.barHover.EbitenImage()
	}
	r.DrawImage(bar, op)

}

// UIScrollViewImpl ...
type UIScrollViewImpl struct {
	Base
	list      *listView
	listRect  image.Rectangle
	scrollBar *scrollBar
	// stroke   *input.Stroke
}

// NewUIScrollViewByList ...
func NewUIScrollViewByList(parent interfaces.Layer, colNames []interface{}, data [][]interface{}, w, h, rowh int, pos *ebitest.Point) interfaces.UIScrollView {
	// img := ebitest.CreateBorderedRectImage(w, h, color.RGBA{0, 0, 0, 255})

	// スクロールビュー全体のベース画像
	eimg := ebiten.NewImage(w, h)
	label := fmt.Sprintf("scroll-%s", utils.RandomLC(8))
	con := &Base{
		label:          label,
		image:          eimg,
		position:       pos,
		layer:          parent,
		hasHoverAction: true,
	}

	// リスト領域
	lw := w - sbWidth
	lh := h - rowh // 列名分の高さをマイナスする

	// リストを作成（ポジションは列名分下げた位置）
	l := newListView(colNames, data, lw, rowh, ebitest.NewPoint(0, float64(rowh)))

	// // スクロールパーツ作成
	// scrollbaseimg := ebitest.CreateRectImage(sbWidth, 10, color.RGBA{255, 255, 255, 255})
	// scrollbarimg := ebitest.CreateRectImage(sbWidth-5, 10, color.RGBA{192, 192, 192, 255})
	// scrollbarhilightimg := ebitest.CreateRectImage(sbWidth-5, 10, color.RGBA{127, 127, 127, 255})

	bgSize := ebitest.NewSize(l.image.Size())
	wscale := float64(lw) / float64(bgSize.W())
	hscale := float64(lh*lh-6) / (float64(bgSize.H()) * wscale)
	// log.Printf("scrollbarScale: %0.1f", hscale)

	sbpos := ebitest.NewPoint(float64(lw), float64(rowh))
	sb := &scrollBar{
		base:        models.NewEbiObject(label, ebiten.NewImageFromImage(scrollbaseimg), con.image, ebitest.NewScale(1.0, float64(lh)/sbHeight), sbpos, 0, true, false, false),
		bar:         models.NewEbiObject(label, ebiten.NewImageFromImage(scrollbarimg), con.image, ebitest.NewScale(1.0, hscale/sbHeight), sbpos, 0, true, false, false),
		barHover:    models.NewEbiObject(label, ebiten.NewImageFromImage(scrollbarhilightimg), con.image, ebitest.NewScale(1.0, hscale/sbHeight), sbpos, 0, true, false, false),
		draggingPos: ebitest.NewPoint(0.0, 0.0),
		scrollMin:   0,
		scrollMax:   float64(lh) - (float64(bgSize.H()) * wscale),
		// scrollbarScale: hscale,
	}

	s := &UIScrollViewImpl{
		Base:      *con,
		list:      l,
		listRect:  image.Rect(0, 0, lw, lh),
		scrollBar: sb,
	}
	return s
}

// Update ...
func (o *UIScrollViewImpl) Update() error {
	// ホイールイベント
	_, dy := ebiten.Wheel()
	// log.Printf("%0.1f < %0.1f && %0.1f > %0.1f", c.contentOffset.y, c.scrollMin, c.contentOffset.y, c.scrollMax)
	if dy < 0 {

		if o.listRect.Min.Y > 0 {
			o.listRect = image.Rect(o.listRect.Min.X, o.listRect.Min.Y+int(dy*2), o.listRect.Max.X, o.listRect.Max.Y+int(dy*2))
			// log.Printf("listRect: %#v", c.listRect)
		}

	} else {

		// 表示領域高さ
		h := float64(o.listRect.Max.Y - o.listRect.Min.Y)
		// コンテンツサイズ
		listSize := ebitest.NewSize(o.list.image.Size())
		// スクロール最大位置
		scrollMax := int((float64(listSize.H()) * o.list.Scale(enum.TypeGlobal).Y()) - h)
		// log.Printf("c.listRect.Min.Y: %d, scrollMax: %d", c.listRect.Min.Y, scrollMax)
		miny := o.listRect.Min.Y + int(dy*2)
		maxy := o.listRect.Max.Y + int(dy*2)
		// スクロールしすぎ防止処理
		if miny > scrollMax {
			y := miny - scrollMax
			miny -= y
			maxy -= y
		} else if miny < 0 {
			y := miny
			miny = 0
			maxy -= y
		}
		o.listRect = image.Rect(o.listRect.Min.X, miny, o.listRect.Max.X, maxy)
	}

	// ドラッグイベント
	if stroke := o.scrollBar.stroke; stroke != nil {
		o.scrollBar.updateStroke(stroke)
		if stroke.IsReleased() {
			o.scrollBar.stroke = nil
			o.listRect = image.Rect(o.listRect.Min.X, o.listRect.Min.Y+int(o.scrollBar.draggingPos.Y()), o.listRect.Max.X, o.listRect.Max.Y+int(o.scrollBar.draggingPos.Y()))
			log.Printf("drag end")
		}
	}

	if o.scrollBar.base.In(ebiten.CursorPosition()) {
		// log.Printf("hover")
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			stroke := input.NewStroke(&input.MouseStrokeSource{})
			// レイヤ内のドラッグ対象のオブジェクトを取得する仕組みが必要
			o.scrollBar.stroke = stroke
			log.Printf("drag start")
		}
	}

	// c.listRect = image.Rect(c.listRect.Min.X, c.listRect.Min.Y+int(dy*5), c.listRect.Max.X, c.listRect.Max.Y+int(dy*5))
	return nil
}

// Draw ...
func (o *UIScrollViewImpl) Draw(screen *ebiten.Image) {
	var op *ebiten.DrawImageOptions

	// 列名を描画
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(o.Position(enum.TypeGlobal).Get())
	screen.DrawImage(o.list.names.DrawListRow(false), op)

	//スクロールバードラッグ位置を反映
	rect := o.listRect
	if o.scrollBar.stroke != nil {
		listSize := ebitest.NewSize(o.list.image.Size())
		dy := int(o.scrollBar.draggingPos.Y())
		if o.listRect.Min.Y+dy < 0 {
			dy = 0 - o.listRect.Min.Y
		} else if o.listRect.Max.Y+dy > listSize.H() {
			dy = listSize.H() - o.listRect.Max.Y
		}
		// log.Printf("dy: %d", dy)
		rect = image.Rect(o.listRect.Min.X, o.listRect.Min.Y+dy, o.listRect.Max.X, o.listRect.Max.Y+dy)
	}

	// リスト部分を作成
	list := o.list.DrawList(rect)

	// リストの表示部分を描画
	op = &ebiten.DrawImageOptions{}
	// 描画位置指定
	op.GeoM.Reset()
	op.GeoM.Scale(o.Scale(enum.TypeGlobal).Get())
	// 描画位置
	op.GeoM.Translate(o.list.Position(enum.TypeGlobal).Get())
	screen.DrawImage(list.SubImage(rect).(*ebiten.Image), op)

	// スクロールバー
	listSize := ebitest.NewSize(list.Size())
	bgScale := o.Scale(enum.TypeGlobal)
	o.scrollBar.Draw(screen, float64(listSize.H()), float64(rect.Min.Y), bgScale.Y())

}

// getFocusedRow カーソルのある行のインデックスと参照を返します
func (o *UIScrollViewImpl) getFocusedRow() (int, *listRowView) {
	drawRect := o.listRect
	l := o.list

	// カーソルIN判定
	isHover := false
	cx, cy := ebiten.CursorPosition()
	bgPos := o.Position(enum.TypeGlobal)
	viewSize := ebitest.NewSize(drawRect.Dx(), drawRect.Dy())
	_, dy := 0, 0
	if int(bgPos.X()) <= cx && int(bgPos.X())+viewSize.W() >= cx {
		// 横位置がスクロールビュー範囲内
		if int(bgPos.Y()) <= cy && int(bgPos.Y())+viewSize.H() >= cy {
			// 縦位置がスクロールビュー範囲内
			// dx = cx - int(bgPos.X())
			dy = cy - int(bgPos.Y())
			isHover = true
		}
	}

	y := 0.0
	top, bottom, min, max := float64(drawRect.Min.Y), float64(drawRect.Max.Y), 0.0, 0.0
	for i, row := range l.rows {
		min = y
		max = y + float64(row.GetRowHeight())

		// 対象行の下端がtop以下あるいは対象行の上端がbottom以上、以外が描画対象
		if !(max <= top || min >= bottom) {
			if isHover {
				if int(min-top) <= dy && int(max-top) >= dy {
					// log.Printf("カーソルは%d行目の範囲内: x=%d, y=%d", i, dx, dy)
					return i, &row
				}
			}
		}
		y += float64(row.GetRowHeight()) + margin
	}
	return -1, nil
}

// GetIndexOfFocusedRow カーソルのある行のインデックスと参照を返します
func (o *UIScrollViewImpl) GetIndexOfFocusedRow() int {
	i, _ := o.getFocusedRow()
	return i
}
