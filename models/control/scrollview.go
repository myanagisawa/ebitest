package control

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/app/g"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/functions"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/models/char"
	"github.com/myanagisawa/ebitest/utils"
)

var (
	marginX, marginY = 2, 3
)

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
	scrollbarBar  *scrollbarBar
	fontSet       *char.Resource
	onHeaderClick func(row interface{}, pos *g.Point, params map[string]interface{})
	onRowClick    func(index int, row interface{}, pos *g.Point, params map[string]interface{})
}

// Children ...
func (o *UIScrollView) Children() []interfaces.UIControl {
	ret := []interfaces.UIControl{}
	ret = append(ret, o.header)
	ret = append(ret, o.listView)
	ret = append(ret, o.listView.Children()...)
	ret = append(ret, o.scrollbarBase)
	ret = append(ret, o.scrollbarBar)
	// log.Printf("*UIScrollView.Children: l=%s, o=%T", o.Label(), o)
	return ret
}

// SetDataSource ...
func (o *UIScrollView) SetDataSource(colNames []interface{}, data [][]interface{}) {
	// log.Printf("fontSet: %#v", o.fontSet)
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
	// log.Printf("listWidth: %d", listWidth)
	// リスト幅から各カラムのマージン分のサイズをマイナス
	calcw := listWidth - (marginX * (len(ratios) + 1))
	// カラムサイズリストを取得
	colWidth := make([]int, len(ratios))
	for i := range ratios {
		colWidth[i] = int(float64(calcw) * ratios[i])
	}
	// log.Printf("colWidth: %#v", colWidth)

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
	row := newListRow("header", o, headers, 0, g.NewPoint(0, 0), listWidth, headerheight)
	o.header = row

	headerheight += marginY
	// スクロール部分の初期化
	list := newListView(fmt.Sprintf("%s.list", o.label), o, g.NewPoint(0, float64(headerheight)))
	o.listView = list

	// 行オブジェクト作成
	totalHeight := 0
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
		row := newListRow(fmt.Sprintf("row-%d", i), o, cols, i, g.NewPoint(0, float64(totalHeight)), listWidth, rowheight)
		o.listView.SetRow(row)
		// height更新
		totalHeight += rowheight + marginY
	}
	// 表示領域確定
	o.listView.setDisplayIndex()

	// スクロールバー部分の初期化
	barBase := newScrollbarBase(fmt.Sprintf("%s.scrollbar.base", o.label), o, g.NewPoint(float64(listWidth-12), float64(headerheight)))
	o.scrollbarBase = barBase

	basePos := barBase.position
	bar := newScrollbarBar(fmt.Sprintf("%s.scrollbar.bar", o.label), o, g.NewPoint(basePos.X()+2, basePos.Y()+3))
	o.scrollbarBar = bar
}

// GetObjects ...
func (o *UIScrollView) GetObjects(x, y int) []interfaces.EbiObject {
	if o.In(x, y) {
		objs := []interfaces.EbiObject{}
		if o.header.In(x, y) {
			objs = append(objs, o.header)
		}
		if o.scrollbarBar.In(x, y) {
			objs = append(objs, o.scrollbarBar)
		}
		if o.scrollbarBase.In(x, y) {
			objs = append(objs, o.scrollbarBase)
		}
		for i := o.listView.displayFrom; i <= o.listView.displayTo; i++ {
			c := o.listView.rows[i]
			if c.In(x, y) {
				objs = append(objs, c)
			}
		}
		if o.listView.In(x, y) {
			objs = append(objs, o.listView)
		}
		objs = append(objs, o)
		// log.Printf("UIScrollView::GetObjects: %#v", objs)
		return objs
	}
	return nil
}

// Update ...
func (o *UIScrollView) Update() error {
	o.listView.Update()
	o.scrollbarBase.Update()
	o.scrollbarBar.Update()

	return nil
}

// Draw ...
func (o *UIScrollView) Draw(screen *ebiten.Image) {
	var op *ebiten.DrawImageOptions

	o.Base.Draw(screen)
	o.listView.Draw(screen)
	o.scrollbarBase.Draw(screen)
	o.scrollbarBar.Draw(screen)

	// ヘッダ描画
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(o.Position(enum.TypeGlobal).Get())
	screen.DrawImage(o.header.image, op)

	// log.Printf("UIScrollView.Draw: l=%s, o=%T", o.Label(), o)
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

// SetRowClickFunc ヘッダクリック、行クリック時の処理を設定します
func (o *UIScrollView) SetRowClickFunc(headerfunc func(row interface{}, pos *g.Point, params map[string]interface{}), rowfunc func(idx int, row interface{}, pos *g.Point, params map[string]interface{})) {
	if headerfunc == nil {
		log.Printf("ヘッダクリック処理の設定をスキップ")
	} else {
		if o.header == nil {
			log.Printf("ヘッダ未定義のためクリック処理の設定をスキップ")
		} else {
			o.onHeaderClick = headerfunc
		}
	}

	if rowfunc == nil {
		log.Printf("行クリック処理の設定をスキップ")
	} else {
		if o.listView == nil {
			log.Printf("リスト未定義のためクリック処理の設定をスキップ")
		} else {
			o.onRowClick = rowfunc
		}
	}
}

// NewUIScrollView ...
func NewUIScrollView(l interfaces.Layer, label string, pos *g.Point, size *g.Size) interfaces.UIScrollView {
	eimg := ebiten.NewImage(size.Get())
	cb := NewControlBase(l, label, eimg, pos, g.DefScale(), true).(*Base)
	o := &UIScrollView{
		Base:    *cb,
		fontSet: char.Res.Get(12, enum.FontStyleGenShinGothicNormal),
	}
	o.onHeaderClick = func(row interface{}, pos *g.Point, params map[string]interface{}) {
		log.Printf("デフォルトヘッダクリックイベントだよ")
	}
	o.onRowClick = func(index int, row interface{}, pos *g.Point, params map[string]interface{}) {
		log.Printf("デフォルト行クリックイベントだよ(%d)", index)
	}
	o.eventHandler.AddEventListener(enum.EventTypeFocus, functions.CommonEventCallback)

	return o
}

/*****************************************************************/

// scrollViewParts パーツの基底
type scrollViewParts struct {
	Base
	parent *UIScrollView
}

// Manager ...
func (o *scrollViewParts) Manager() interfaces.GameManager {
	return o.parent.Manager()
}

func (o *scrollViewParts) Layer() interfaces.Layer {
	return o.parent.Layer()
}

// Position ...
func (o *scrollViewParts) Position(t enum.ValueTypeEnum) *g.Point {
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
	return g.NewPoint(gx, gy)
}

// In ...
func (o *scrollViewParts) In(x, y int) bool {
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
	// return controlIn(x, y,
	// 	o.Position(enum.TypeGlobal),
	// 	g.NewSize(o.image.Size()),
	// 	o.Scale(enum.TypeGlobal),
	// 	o.Layer().Frame().Position(enum.TypeGlobal),
	// 	o.Layer().Frame().Size())
}

// Update ...
func (o *scrollViewParts) Update() error {
	// if o.hasHoverAction {
	// 	o.hover = o.In(ebiten.CursorPosition())
	// 	if o.hover {
	// 		log.Printf("hover: %s", o.label)
	// 	}
	// }
	return nil
}

// Draw ...
func (o *scrollViewParts) Draw(screen *ebiten.Image) {
	var op *ebiten.DrawImageOptions

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(o.Position(enum.TypeGlobal).Get())

	screen.DrawImage(o.image, op)

	// log.Printf("scrollViewParts.Draw: l=%s, o=%T", o.Label(), o)
}

func newScrollViewParts(label string, parent *UIScrollView, eimg *ebiten.Image, pos *g.Point) *scrollViewParts {
	if eimg == nil {
		eimg = ebiten.NewImage(0, 0)
	}
	cb := NewControlBase(parent.layer, label, eimg, pos, g.DefScale(), true).(*Base)

	o := &scrollViewParts{
		Base:   *cb,
		parent: parent,
	}
	return o
}

/*****************************************************************/

// listView ...
type listView struct {
	scrollViewParts
	scrollingPos *g.Point
	rows         []*listRow
	size         *g.Size
	displayFrom  int
	displayTo    int
}

// Children ...
func (o *listView) Children() []interfaces.UIControl {
	ret := make([]interfaces.UIControl, len(o.rows))
	for i, row := range o.rows {
		ret[i] = row
	}
	// log.Printf("*listView.Children: l=%s, o=%T", o.Label(), o)
	return ret
}

func (o *listView) calcScrollingPos(dy int) *g.Point {
	sx, sy := o.scrollingPos.GetInt()
	if dy != 0 {
		sy += dy
		if sy < 0 {
			// 上に余白ができる
			sy = 0
		} else {
			ih := o.size.H()
			_, ph := o.parent.listViewSize()
			if sy+ph > ih {
				// 下に余白ができる
				sy = ih - ph
			}
		}
	}
	// log.Printf("calcScrollingPos: %d", sy)
	return g.NewPoint(float64(sx), float64(sy))
}

// ScrollingPos ...
func (o *listView) ScrollingPos() *g.Point {
	if o.moving == nil {
		return o.scrollingPos
	}
	// スクロールバードラッグ中
	_, dy := o.moving.GetInt()
	return o.calcScrollingPos(dy)
}

func (o *listView) DidWheel(dx, dy float64) {
	if dy != 0 {
		o.scrollingPos.Set(o.calcScrollingPos(int(dy * 2)).Get())
		// 表示領域確定
		o.setDisplayIndex()
	}
}

// Update ...
func (o *listView) Update() error {
	o.scrollViewParts.Update()
	// ホイールイベント
	// _, dy := ebiten.Wheel()
	// if dy != 0 {
	// 	o.scrollingPos.Set(o.calcScrollingPos(int(dy * 2)).Get())
	// 	// 表示領域確定
	// 	o.setDisplayIndex()
	// }

	// 表示対象業のアップデート処理
	for i := o.displayFrom; i <= o.displayTo; i++ {
		o.rows[i].Update()
	}
	return nil
}

// setDisplayIndex 表示対象行の設定を行います. 表示対象領域の変化があった場合に再計算させる処理なので、スクロール処理に関連づけて実行する
func (o *listView) setDisplayIndex() {
	// スクロール量
	_, sy := o.ScrollingPos().GetInt()
	_, lh := o.parent.listViewSize()
	topY, bottomY := 0, lh

	from, to := -1, -1
	for i := range o.rows {
		row := o.rows[i]
		_, rh := row.image.Size()

		ty := int(row.position.Y()) - sy

		// 描画領域判定
		if ty <= topY {
			if ty+rh > topY {
				// 上端に一部隠れた状態=表示対象の先頭
				from = i
			}
		} else if ty+rh > bottomY {
			to = i
			break
		} else {
			// 通常描画領域
			if from == -1 {
				from = i
			}
		}
	}
	if to == -1 {
		to = len(o.rows) - 1
	}
	o.displayFrom = from
	o.displayTo = to
}

// Draw ...
func (o *listView) Draw(screen *ebiten.Image) {
	var op *ebiten.DrawImageOptions

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(o.Position(enum.TypeGlobal).Get())

	// スクロール量
	_, sy := o.ScrollingPos().GetInt()

	// リストの描画位置
	x, y := o.Position(enum.TypeGlobal).GetInt()
	// リスト上下端位置
	topY := y
	_, lh := o.parent.listViewSize()
	bottomY := y + lh
	for i := o.displayFrom; i <= o.displayTo; i++ {
		row := o.rows[i]
		_, ry := row.position.GetInt()
		ty := y + ry - sy

		// 行の描画
		op = &ebiten.DrawImageOptions{}
		r, g, b, a := 1.0, 1.0, 1.0, 1.0
		if row.hover {
			// log.Printf("ホバー行: %s(%d)", row.label, row.index)
			r, g, b, a = 0.75, 0.75, 0.75, 1.0
		}
		op.ColorM.Scale(r, g, b, a)

		if i == o.displayFrom {
			// 先頭データ
			a := topY - ty
			if a < 0 {
				// 行の間のマージン部分が上端にかかっている場合にマイナスで出る
				// そのままだと先頭セルが上にズレるので0で上書き
				a = 0
			}
			op.GeoM.Translate(float64(x), float64(ty+a))

			// 画像の表示範囲を計算
			rw, rh := row.image.Size()
			fr := image.Rect(0, a, rw, rh+a)
			screen.DrawImage(row.image.SubImage(fr).(*ebiten.Image), op)
		} else if i == o.displayTo {
			// 末尾データ
			a := ty - bottomY
			op.GeoM.Translate(float64(x), float64(ty))

			// 画像の表示範囲を計算
			rw, _ := row.image.Size()
			fr := image.Rect(0, 0, rw, -a)
			screen.DrawImage(row.image.SubImage(fr).(*ebiten.Image), op)
		} else {
			// 境界以外
			op.GeoM.Translate(float64(x), float64(ty))
			screen.DrawImage(row.image, op)
		}
	}
	// log.Printf("listView.Draw: l=%s, o=%T", o.Label(), o)

}

// SetRow ...
func (o *listView) SetRow(row *listRow) {
	_, rh := row.image.Size()
	o.rows = append(o.rows, row)
	o.size = g.NewSize(o.size.W(), o.size.H()+rh+marginY)
}

func newListView(label string, parent *UIScrollView, pos *g.Point) *listView {
	pw, ph := parent.image.Size()

	eimg := ebiten.NewImage(pw, ph-int(pos.Y()))

	sb := newScrollViewParts(label, parent, eimg, pos)

	o := &listView{
		scrollViewParts: *sb,
		scrollingPos:    g.NewPoint(0, 0),
		size:            g.NewSize(pw, -marginY),
	}
	o.eventHandler.AddEventListener(enum.EventTypeWheel, functions.CommonEventCallback)
	return o
}

/*****************************************************************/

// listRow ...
type listRow struct {
	scrollViewParts
	index  int
	source []string
}

// IsHeader ...
func (o *listRow) IsHeader() bool {
	return (o.label == "header")
}

// Index ...
func (o *listRow) Index() int {
	return o.index
}

// Source ...
func (o *listRow) Source() []string {
	return o.source
}

// Parent ...
func (o *listRow) Parent() interfaces.UIScrollView {
	return o.parent
}

// Position ...
func (o *listRow) Position(t enum.ValueTypeEnum) *g.Point {
	// スクロールバー位置: x = リスト位置(sy)*スクロール領域サイズ(sh) / リストサイズ(lh)
	by := 0.0
	if !o.IsHeader() {
		_, sy := o.parent.listView.ScrollingPos().Get()
		py := o.position.Y()
		by = py - sy
	} else {
		by = o.position.Y()
	}
	if t == enum.TypeLocal {
		return g.NewPoint(0, by)
	}
	gy := 0.0
	if !o.IsHeader() {
		_, gy = o.parent.listView.Position(enum.TypeGlobal).Get()
	} else {
		_, gy = o.parent.Position(enum.TypeGlobal).Get()
	}
	_, sy := o.Scale(enum.TypeGlobal).Get()
	gy += by * sy
	return g.NewPoint(0, gy)
}

// In ...
func (o *listRow) In(x, y int) bool {
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
	// return controlIn(x, y,
	// 	o.Position(enum.TypeGlobal),
	// 	g.NewSize(o.image.Size()),
	// 	o.Scale(enum.TypeGlobal),
	// 	o.Layer().Frame().Position(enum.TypeGlobal),
	// 	o.Layer().Frame().Size())
}

// Update ...
func (o *listRow) Update() error {
	// o.hover = o.In(ebiten.CursorPosition())
	return nil
}

func newListRow(label string, parent *UIScrollView, columns []*column, index int, pos *g.Point, width, height int) *listRow {
	// 行画像を作成
	img := utils.CreateRectImage(width, height, &color.RGBA{0, 0, 0, 32}).(draw.Image)

	sl := make([]string, len(columns))
	cx := marginX
	for i := range columns {
		col := columns[i]
		columnImageBase := utils.CreateRectImage(col.width, height, &color.RGBA{127, 127, 127, 64}).(draw.Image)

		// データタイプごとの描画
		switch col.ds.(type) {
		case image.Image:
			// 画像

			// ソース文字列
			sl[i] = fmt.Sprintf("%s", "image column")
		case int:
			// テキスト（数値）
			tx := col.width - col.padding[1]
			for j := len(col.sources) - 1; j >= 0; j-- {
				t := col.sources[j]
				tx -= t.Bounds().Size().X
				// log.Printf("col: %d, j: %d, width: %d, pad: %d, tx: %d", col.colIndex, j, col.width, col.padding[1], tx)
				columnImageBase = utils.StackImage(columnImageBase, t, image.Point{tx, col.padding[0]})
			}
			for j := range col.sources {
				t := col.sources[j]
				columnImageBase = utils.StackImage(columnImageBase, t, image.Point{tx, col.padding[0]})
				tx += t.Bounds().Size().X
			}
			// ソース文字列
			sl[i] = fmt.Sprintf("%d", col.ds)
		case string:
			tx := col.padding[3]
			for j := range col.sources {
				t := col.sources[j]
				columnImageBase = utils.StackImage(columnImageBase, t, image.Point{tx, col.padding[0]})
				tx += t.Bounds().Size().X
			}
			// ソース文字列
			sl[i] = fmt.Sprintf("%s", col.ds)
		default:
			panic("invalid type")
		}

		// カラム画像を行画像上に描画
		img = utils.StackImage(img, columnImageBase, image.Point{cx, 0})
		cx += columnImageBase.Bounds().Size().X + marginX
	}

	sb := newScrollViewParts(label, parent, ebiten.NewImageFromImage(img), pos)
	o := &listRow{
		scrollViewParts: *sb,
		index:           index,
		source:          sl,
	}
	o.eventHandler.AddEventListener(enum.EventTypeFocus, functions.CommonEventCallback)
	o.eventHandler.AddEventListener(enum.EventTypeClick, func(ev interfaces.EventOwner, pos *g.Point, params map[string]interface{}) {
		if row, ok := ev.(interfaces.ListRow); ok {
			if p, ok := row.Parent().(*UIScrollView); ok {
				if row.IsHeader() {
					p.onHeaderClick(row, pos, params)
				} else {
					p.onRowClick(row.Index(), row, pos, params)
				}

				tname := fmt.Sprintf("%s", reflect.TypeOf(o))
				log.Printf("スクロールビューのクリック(%d): %s", row.Index(), tname)
			}
		}
	})

	return o
}

/*****************************************************************/

// scrollbarBase ...
type scrollbarBase struct {
	scrollViewParts
}

func newScrollbarBase(label string, parent *UIScrollView, pos *g.Point) *scrollbarBase {
	_, ph := parent.listViewSize()
	scrollbaseimg := utils.CreateRectImage(12, ph, &color.RGBA{223, 223, 223, 255})
	eimg := ebiten.NewImageFromImage(scrollbaseimg)

	sb := newScrollViewParts(label, parent, eimg, pos)
	o := &scrollbarBase{
		scrollViewParts: *sb,
	}
	return o
}

/*****************************************************************/

// scrollbarBar ...
type scrollbarBar struct {
	scrollViewParts
}

func (o *scrollbarBar) DidStroke(dx, dy float64) {
	o.parent.listView.SetMoving(0, dy)
	o.parent.listView.setDisplayIndex()
}

func (o *scrollbarBar) FinishStroke() {
	_, dy := o.parent.listView.moving.GetInt()
	o.parent.listView.scrollingPos.Set(o.parent.listView.calcScrollingPos(dy).Get())
	o.parent.listView.moving = nil
}

// Position ...
func (o *scrollbarBar) Position(t enum.ValueTypeEnum) *g.Point {
	// スクロールバー位置: x = リスト位置(sy)*スクロール領域サイズ(sh) / リストサイズ(lh)
	by := 0.0
	{
		_, sy := o.parent.listView.ScrollingPos().GetInt()
		_, sh := o.parent.scrollbarBase.image.Size()
		sh -= 6 // ベース領域から、バーのマージン上下各3px分を引く
		lh := o.parent.listView.size.H()
		by = math.Abs(float64(sy)) * float64(sh) / float64(lh)
	}

	if t == enum.TypeLocal {
		return g.NewPoint(o.position.X(), o.position.Y()+by)
	}
	gx, gy := o.parent.Position(enum.TypeGlobal).Get()
	sx, sy := o.Scale(enum.TypeGlobal).Get()
	gx += o.position.X() * sx
	gy += (o.position.Y() + by) * sy
	return g.NewPoint(gx, gy)
}

// In ...
func (o *scrollbarBar) In(x, y int) bool {
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
	// return controlIn(x, y,
	// 	o.Position(enum.TypeGlobal),
	// 	g.NewSize(o.image.Size()),
	// 	o.Scale(enum.TypeGlobal),
	// 	o.Layer().Frame().Position(enum.TypeGlobal),
	// 	o.Layer().Frame().Size())
}

// Update ...
func (o *scrollbarBar) Update() error {
	// o.hover = o.In(ebiten.CursorPosition())
	return nil
}

// Draw ...
func (o *scrollbarBar) Draw(screen *ebiten.Image) {
	var op *ebiten.DrawImageOptions

	op = &ebiten.DrawImageOptions{}

	r, g, b, a := 1.0, 1.0, 1.0, 1.0
	if o.hover {
		r, g, b, a = 0.75, 0.75, 0.75, 1.0
	}
	op.ColorM.Scale(r, g, b, a)

	// x, y := o.Position(enum.TypeGlobal).Get()
	op.GeoM.Translate(o.Position(enum.TypeGlobal).Get())
	screen.DrawImage(o.image, op)
	// log.Printf("scrollbarBar.Draw: l=%s, o=%T", o.Label(), o)
}

func newScrollbarBar(label string, parent *UIScrollView, pos *g.Point) *scrollbarBar {
	// リスト表示領域高さ
	_, ph := parent.listViewSize()
	// リスト全体の高さ
	lh := parent.listView.size.H()

	// barの高さ
	barheight := int(float64(ph*ph) / float64(lh))
	barheight -= 3

	scrollbarimg := utils.CreateRectImage(8, barheight, &color.RGBA{192, 192, 192, 255})
	eimg := ebiten.NewImageFromImage(scrollbarimg)

	sb := newScrollViewParts(label, parent, eimg, pos)
	o := &scrollbarBar{
		scrollViewParts: *sb,
	}
	o.eventHandler.AddEventListener(enum.EventTypeFocus, functions.CommonEventCallback)
	o.eventHandler.AddEventListener(enum.EventTypeDragging, functions.CommonEventCallback)
	o.eventHandler.AddEventListener(enum.EventTypeDragDrop, functions.CommonEventCallback)

	return o

}

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
func (o *column) GetSourcesSize() *g.Size {
	w, h := 0, 0
	for i := range o.sources {
		s := o.sources[i]
		sp := s.Bounds().Size()
		w += sp.X
		if h < sp.Y {
			h = sp.Y
		}
	}
	return g.NewSize(w, h+o.padding[0]+o.padding[2])
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
