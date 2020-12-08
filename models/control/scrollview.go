package control

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/ebitest"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/models/char"
	"github.com/myanagisawa/ebitest/utils"
)

var (
	marginX, marginY = 2, 3
)

// // GetSourceImages ...
// func GetSourceImages(fontSet *char.Resource, source interface{}) []image.Image {
// 	ret := []image.Image{}

// 	switch val := source.(type) {
// 	case image.Image:
// 		ret = append(ret, val)
// 	case int:
// 		ret = fontSet.GetByString(fmt.Sprintf("%d", val))
// 	case string:
// 		ret = fontSet.GetByString(val)
// 	default:
// 		panic("invalid type")
// 	}
// 	return ret
// }

// GetColumnWidthRatios カラムごとの幅の比率を取得します
func GetColumnWidthRatios(columns *columnSet) []float64 {

	maxWidths := make([]int, columns.colCount)
	for i := range columns.columns {
		col := columns.columns[i]
		if col.GetSourcesSize().W() > maxWidths[col.colIndex] {
			maxWidths[col.colIndex] = col.GetSourcesSize().W()
		}
	}
	// 列名行も対象
	for i := range columns.header {
		col := columns.header[i]
		if col.GetSourcesSize().W() > maxWidths[col.colIndex] {
			maxWidths[col.colIndex] = col.GetSourcesSize().W()
		}
	}

	// 最大幅での各列のサイズ比を計算
	totalWidth := 0.0
	ratio := make([]float64, columns.colCount)
	for i := range maxWidths {
		totalWidth += float64(maxWidths[i])
	}
	for i := range maxWidths {
		ratio[i] = float64(maxWidths[i]) / totalWidth
	}
	return ratio
}

/*****************************************************************/

// UIScrollView ...
type UIScrollView struct {
	Base
	header        *listRow
	listView      *listView
	scrollbarBase *scrollbarBase
	fontSet       *char.Resource
}

// SetDataSource ...
func (o *UIScrollView) SetDataSource(colNames []interface{}, data [][]interface{}) {
	// カラムデータセットをまずは作ってしまう
	columns := &columnSet{}
	for i := range data {
		row := data[i]
		for j := range row {
			columns.AddColumn(newColumn("", o, i, j, row[j]))
		}
	}
	for i := range colNames {
		columns.AddHeader(newColumn("", o, 0, i, colNames[i]))
	}

	// カラム幅の計算
	ratios := GetColumnWidthRatios(columns)
	listWidth, _ := o.image.Size()
	log.Printf("listWidth: %d", listWidth)
	// リスト幅から各カラムのマージン分のサイズをマイナス
	calcw := listWidth - (marginX * (len(ratios) + 1))
	// カラムサイズリストを取得
	colWidth := make([]int, len(ratios))
	for i := range ratios {
		colWidth[i] = int(float64(calcw) * ratios[i])
	}
	log.Printf("colWidth: %#v", colWidth)

	// ヘッダ作成
	headers := columns.GetHeader()
	// カラムサイズを設定
	headerheight := 0
	for j := range headers {
		headers[j].width = colWidth[j]
		s := headers[j].GetSourcesSize()
		if headerheight < s.H() {
			headerheight = s.H()
		}
	}
	row := newListRow("header", o, headers, 0, listWidth, headerheight)
	o.header = row

	headerheight += marginY
	// スクロール部分の初期化
	list := newListView(fmt.Sprintf("%s.list", o.label), o, ebitest.NewPoint(0, float64(headerheight)))
	o.listView = list

	// 行オブジェクト作成
	// totalHeight := marginY
	for i := 0; i < columns.rowCount; i++ {
		// 対象行のカラムリスト取得
		cols := columns.GetByRowIndex(i)
		// カラムサイズを設定
		rowheight := 0
		for j := range cols {
			cols[j].width = colWidth[j]
			s := cols[j].GetSourcesSize()
			if rowheight < s.H() {
				rowheight = s.H()
			}
		}
		// 行を作成
		row := newListRow(fmt.Sprintf("row-%d", i), o, cols, i, listWidth, rowheight)
		o.listView.SetRow(row)
		// totalHeight += rowheight + marginY
	}

	// スクロールバー部分の初期化
	barBase := newScrollbarBase(fmt.Sprintf("%s.scrollbar.base", o.label), o, ebitest.NewPoint(float64(listWidth-15), float64(headerheight)))
	o.scrollbarBase = barBase

}

// Update ...
func (o *UIScrollView) Update() error {
	o.listView.Update()
	o.scrollbarBase.Update()

	return nil
}

// Draw ...
func (o *UIScrollView) Draw(screen *ebiten.Image) {
	var op *ebiten.DrawImageOptions

	o.Base.Draw(screen)
	o.listView.Draw(screen)
	o.scrollbarBase.Draw(screen)

	// ヘッダ描画
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(o.Position(enum.TypeGlobal).Get())
	screen.DrawImage(o.header.image, op)
}

// listViewSize スクロールビューの中のリスト領域のサイズを返す
func (o *UIScrollView) listViewSize() (int, int) {
	// スクロールビューサイズ
	iw, ih := o.image.Size()
	// ヘッダサイズ
	_, hh := o.header.image.Size()

	//リスト表示領域の調整をここで実施
	w, h := iw, ih-(hh+marginY)
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
	size         *ebitest.Size
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
		ih := o.size.H()
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

	// スクロール量
	_, sy := o.scrollingPos.GetInt()
	// リストの描画位置
	x, y := o.Position(enum.TypeGlobal).GetInt()
	// リスト上下端位置
	topY := y
	_, lh := o.parent.listViewSize()
	bottomY := y + lh
	th := 0
	for i := range o.rows {
		row := o.rows[i]
		rw, rh := row.image.Size()

		ty := y + th - sy
		// 描画領域判定
		if ty <= topY {
			if ty+rh > topY {
				// 上端に一部隠れた状態
				a := topY - ty

				op = &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x), float64(ty+a))

				fr := image.Rect(0, a, rw, rh+a)
				screen.DrawImage(row.image.SubImage(fr).(*ebiten.Image), op)
			}
		} else if ty+rh > bottomY {
			if ty <= bottomY {
				// 下端に一部隠れた状態
				a := ty - bottomY

				op = &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x), float64(ty))

				fr := image.Rect(0, 0, rw, -a)
				screen.DrawImage(row.image.SubImage(fr).(*ebiten.Image), op)
			} else {
				// ここに到達したらもう以降は表示範囲外なのでループを抜ける
				break
			}
		} else {
			// 通常描画領域
			op = &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(ty))

			screen.DrawImage(row.image, op)
		}
		th += rh + marginY
	}

}

// SetRow ...
func (o *listView) SetRow(row *listRow) {
	_, rh := row.image.Size()
	o.rows = append(o.rows, row)
	o.size = ebitest.NewSize(o.size.W(), o.size.H()+rh+marginY)
}

func newListView(label string, parent *UIScrollView, pos *ebitest.Point) *listView {
	// img := ebitest.Images["world"]
	// eimg := ebiten.NewImageFromImage(img)

	// positionは親positionからのdeltaを指定する
	cb := Base{
		label: label,
		// image:          eimg,
		position:       pos,
		scale:          ebitest.NewScale(1.0, 1.0),
		hasHoverAction: false,
	}

	pw, _ := parent.image.Size()
	o := &listView{
		scrollViewParts: scrollViewParts{
			Base:   cb,
			parent: parent,
		},
		scrollingPos: ebitest.NewPoint(0, 0),
		size:         ebitest.NewSize(pw, 0),
	}
	return o
}

/*****************************************************************/

// listRow ...
type listRow struct {
	scrollViewParts
	index int
}

func newListRow(label string, parent *UIScrollView, columns []*column, index, width, height int) *listRow {
	// 行画像を作成
	img := ebitest.CreateRectImage(width, height, &color.RGBA{0, 0, 0, 32}).(draw.Image)

	cx := marginX
	for i := range columns {
		col := columns[i]
		columnImageBase := ebitest.CreateRectImage(col.width, height, &color.RGBA{127, 127, 127, 64}).(draw.Image)

		// データタイプごとの描画
		switch col.ds.(type) {
		case image.Image:
			// 画像
		case int:
			// テキスト（数値）
			tx := col.padding[3]
			for j := range col.sources {
				t := col.sources[j]
				columnImageBase = utils.StackImage(columnImageBase, t, image.Point{tx, col.padding[0]})
				tx += t.Bounds().Size().X
			}

		case string:
			tx := col.padding[3]
			for j := range col.sources {
				t := col.sources[j]
				columnImageBase = utils.StackImage(columnImageBase, t, image.Point{tx, col.padding[0]})
				tx += t.Bounds().Size().X
			}
		default:
			panic("invalid type")
		}

		// カラム画像を行画像上に描画
		img = utils.StackImage(img, columnImageBase, image.Point{cx, 0})
		cx += columnImageBase.Bounds().Size().X + marginX
	}
	cb := Base{
		label:          label,
		scale:          ebitest.NewScale(1.0, 1.0),
		image:          ebiten.NewImageFromImage(img),
		hasHoverAction: true,
	}
	o := &listRow{
		scrollViewParts: scrollViewParts{
			Base:   cb,
			parent: parent,
		},
		index: index,
	}

	return o
}

// func newListRowOld(label string, parent *UIScrollView, idx int, row []interface{}) *listRow {
// 	cb := Base{
// 		label:          label,
// 		scale:          ebitest.NewScale(1.0, 1.0),
// 		hasHoverAction: true,
// 	}

// 	o := &listRow{
// 		scrollViewParts: scrollViewParts{
// 			Base:   cb,
// 			parent: parent,
// 		},
// 		index:  idx,
// 		source: row,
// 	}

// 	w, _ := parent.image.Size()
// 	img := ebitest.CreateRectImage(w, 30, &color.RGBA{127, 127, 127, 127}).(draw.Image)

// 	// テキスト画像生成
// 	splitter := parent.fontSet.GetByString(" | ")
// 	tx := 0
// 	for i := range row {
// 		str := fmt.Sprintf("%v", row[i])
// 		col := parent.fontSet.GetByString(str)
// 		// o.texts = append(o.texts, col...)
// 		for j := range col {
// 			t := col[j]
// 			img = utils.StackImage(img, t, image.Point{tx, 5})
// 			tx += t.Bounds().Size().X
// 		}
// 		// カラム区切りの文字列描画
// 		for j := range splitter {
// 			img = utils.StackImage(img, splitter[j], image.Point{tx, 5})
// 			tx += splitter[j].Bounds().Size().X
// 		}
// 	}

// 	o.image = ebiten.NewImageFromImage(img)

// 	return o
// }

/*****************************************************************/

// columnSet ...
type columnSet struct {
	header   []*column
	columns  []*column
	colCount int
	rowCount int
}

func (o *columnSet) AddHeader(col *column) {
	o.header = append(o.header, col)
}

func (o *columnSet) AddColumn(col *column) {
	o.columns = append(o.columns, col)
	if col.colIndex >= o.colCount {
		o.colCount = col.colIndex + 1
	}
	if col.rowIndex >= o.rowCount {
		o.rowCount = col.rowIndex + 1
	}
}

func (o *columnSet) Get(row, col int) *column {
	for i := range o.columns {
		c := o.columns[i]
		if c.rowIndex == row && c.colIndex == col {
			return c
		}
	}
	return nil
}

func (o *columnSet) GetByRowIndex(idx int) []*column {
	ret := make([]*column, o.colCount)
	for i := 0; i < o.colCount; i++ {
		ret[i] = o.Get(idx, i)
	}
	return ret
}

func (o *columnSet) GetHeader() []*column {
	ret := make([]*column, o.colCount)
	for i := range o.header {
		h := o.header[i]
		ret[h.colIndex] = h
	}
	return ret
}

/*****************************************************************/

// column ...
type column struct {
	rowIndex int
	colIndex int
	sources  []image.Image
	ds       interface{}
	width    int
	padding  []int
	align    string
}

// GetSourcesSize sourcesのサイズを返します
func (o *column) GetSourcesSize() *ebitest.Size {
	w, h := 0, 0
	for i := range o.sources {
		s := o.sources[i]
		sp := s.Bounds().Size()
		w += sp.X
		if h < sp.Y {
			h = sp.Y
		}
	}
	return ebitest.NewSize(w, h+o.padding[0]+o.padding[2])
}

// newColumn columデータを作成します
func newColumn(label string, parent *UIScrollView, rowIdx, colIdx int, c interface{}) *column {
	o := &column{
		rowIndex: rowIdx,
		colIndex: colIdx,
		ds:       c,
	}

	switch val := c.(type) {
	case image.Image:
		o.sources = []image.Image{val}
		o.padding = []int{0, 0, 0, 0}
		o.align = "center"
	case int:
		o.sources = parent.fontSet.GetByString(fmt.Sprintf("%d", val))
		o.padding = []int{3, 3, 3, 3}
		o.align = "right"
	case string:
		o.sources = parent.fontSet.GetByString(val)
		o.padding = []int{3, 3, 3, 3}
		o.align = "left"
	default:
		panic("invalid type")
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
