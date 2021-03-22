package control

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/app/g"
	"github.com/myanagisawa/ebitest/example/t7/lib/interfaces"
	"github.com/myanagisawa/ebitest/example/t7/lib/models/char"
	"github.com/myanagisawa/ebitest/example/t7/lib/models/event"
	"github.com/myanagisawa/ebitest/example/t7/lib/utils"
)

const (
	// LineSpacing 行間の標準サイズ
	LineSpacing = 2.0
)

// CreateImageByDataSource 対象データの画像を返します
func CreateImageByDataSource(ds interface{}, fontSet *char.Resource) image.Image {
	switch val := ds.(type) {
	case image.Image:
		return val
	case int:
		return fontSet.GetStringImage(fmt.Sprintf("%d", val))
	case string:
		return fontSet.GetStringImage(val)
	default:
		panic("invalid type")
	}
}

// UIScrollView ...
type UIScrollView struct {
	UIControl
	header *ScrollViewHeader
	list   *ScrollViewList
	bar    *UIControl
}

// NewUIScrollView ...
func NewUIScrollView(base *UIControl, header *ScrollViewHeader, list *ScrollViewList, bar *UIControl) interfaces.UIControl {
	return &UIScrollView{
		UIControl: *base,
		header:    header,
		list:      list,
		bar:       bar,
	}
}

// SetHeader ...
func (o *UIScrollView) SetHeader(header *ScrollViewHeader) {
	o.header = header
}

// SetList ...
func (o *UIScrollView) SetList(list *ScrollViewList) {
	o.list = list
}

// SetScrollBar ...
func (o *UIScrollView) SetScrollBar(bar *UIControl) {
	o.bar = bar
}

// SetHeaderRow ...
func (o *UIScrollView) SetHeaderRow(dataSet []interface{}) {
	o.header.fontSet = char.Res.Get(12, enum.FontStyleGenShinGothicBold)
	o.header.ds = dataSet

	o.header.dsImages = make([]image.Image, len(dataSet))
	for i, ds := range dataSet {
		o.header.dsImages[i] = CreateImageByDataSource(ds, o.header.fontSet)
	}
	// 行画像のUPDATE
	o.UpdateImages()
}

// AppendRows ...
func (o *UIScrollView) AppendRows(dataSet [][]interface{}) {
	_, listSize := o.list.Bound().ToPosSize()
	for _, dataRow := range dataSet {
		row := NewScrollViewRow(o.scene, dataRow)

		index := len(o.list.children)
		if index > 0 {
			prevRow := o.list.children[index-1]
			row.index = index

			// 追加行のbound情報を更新
			pos := g.NewPoint(row.Bound().Min.X(), prevRow.Bound().Max.Y()+LineSpacing)
			size := g.NewSize(listSize.W(), int(row.Bound().Max.Y()))
			row.Bound().SetPosSize(pos, size)
		}

		o.list.AppendChild(row)
	}

	// 行画像のUPDATE
	o.UpdateImages()
}

// GetColumnWidthRatios カラムごとの幅の比率を取得します
func (o *UIScrollView) GetColumnWidthRatios() []float64 {

	maxWidths := make([]int, len(o.header.dsImages))
	// ヘッダ
	for i, col := range o.header.dsImages {
		size := col.Bounds().Size()
		if size.X > maxWidths[i] {
			maxWidths[i] = size.X
		}
	}
	// リスト
	for _, row := range o.list.children {
		if r, ok := row.(*ScrollViewRow); ok {
			for i, col := range r.dsImages {
				size := col.Bounds().Size()
				if size.X > maxWidths[i] {
					maxWidths[i] = size.X
				}
			}
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

	return ratio
}

// UpdateImages ...
func (o *UIScrollView) UpdateImages() {
	ratios := o.GetColumnWidthRatios()
	// ヘッダ更新
	if o.header != nil {
		_, size := o.header.bound.ToPosSize()
		// ベースの行画像を作成
		img := utils.CreateRectImage(size.W(), size.H(), &color.RGBA{0, 0, 0, 0}).(draw.Image)

		rowW := float64(size.W())
		cx := 0
		for i, ds := range o.header.dsImages {
			colW := int(rowW * ratios[i])
			columnImageBase := utils.CreateRectImage(colW-2, size.H(), o.header.baseColor).(draw.Image)
			// データ画像を重ねる
			columnImageBase = utils.StackImage(columnImageBase, ds, image.Point{3, 3})

			// カラム画像を行画像上に描画
			img = utils.StackImage(img, columnImageBase, image.Point{cx, 0})
			cx += colW
		}
		o.header.drawer = NewDefaultDrawer(ebiten.NewImageFromImage(img))
	}

	// リスト更新
	if o.list != nil {
		for idx, row := range o.list.children {
			if r, ok := row.(*ScrollViewRow); ok {
				_, size := r.bound.ToPosSize()
				// ベースの行画像を作成
				img := utils.CreateRectImage(size.W(), size.H(), &color.RGBA{0, 0, 0, 0}).(draw.Image)

				rowW := float64(size.W())
				cx := 0
				for i, ds := range r.dsImages {
					colW := int(rowW * ratios[i])
					columnImageBase := utils.CreateRectImage(colW-2, size.H(), r.baseColor).(draw.Image)
					// データ画像を重ねる
					columnImageBase = utils.StackImage(columnImageBase, ds, image.Point{3, 3})

					// カラム画像を行画像上に描画
					img = utils.StackImage(img, columnImageBase, image.Point{cx, 0})
					cx += colW
				}
				o.list.children[idx].(*ScrollViewRow).drawer = NewDefaultDrawer(ebiten.NewImageFromImage(img))
			}
		}
	}
	// スクロールバー更新
	if sb := o.list.ScrollBound(true); sb != nil {
		_, ls := o.list.bound.ToPosSize()
		// リスト全体の高さ
		_, ps := sb.ToPosSize()

		//高さを再計算(リスト表示領域の高さ * リスト表示領域の高さ / リスト全体の高さ)
		barheight := int(float64(ls.H()*(ls.H()-(LineSpacing*2))) / float64(ps.H()))
		// barheight -= 3

		// サイズ情報をアップデート
		slider := o.bar.children[0]
		pos, size := slider.Bound().ToPosSize()
		o.bar.children[0].Bound().SetPosSize(pos, g.NewSize(size.W(), barheight))
	}

}

// ScrollViewHeader ...
type ScrollViewHeader struct {
	UIControl
	baseColor *color.RGBA
	fontSet   *char.Resource
	ds        []interface{}
	dsImages  []image.Image
}

// NewScrollViewHeader ...
func NewScrollViewHeader(base *UIControl, color *color.RGBA) interfaces.UIControl {
	return &ScrollViewHeader{
		UIControl: *base,
		baseColor: color,
	}
}

// GetControls ...
func (o *ScrollViewHeader) GetControls() []interfaces.UIControl {
	if o._childrenCache != nil {
		return o._childrenCache
	}
	ret := _getControls(o)
	o._childrenCache = ret
	return ret
}

// ScrollViewList ...
type ScrollViewList struct {
	UIControl
}

// NewScrollViewList ...
func NewScrollViewList(base *UIControl) interfaces.UIControl {
	return &ScrollViewList{
		UIControl: *base,
	}
}

// GetControls ...
func (o *ScrollViewList) GetControls() []interfaces.UIControl {
	if o._childrenCache != nil {
		return o._childrenCache
	}
	ret := _getControls(o)
	o._childrenCache = ret
	return ret
}

// ScrollBound ...
func (o *ScrollViewList) ScrollBound(withoutMoving bool) *g.Bound {
	if children := o.GetChildren(); children != nil {
		fc := children[0].(*ScrollViewRow)
		lc := children[len(children)-1].(*ScrollViewRow)

		min := fc.Bound().Min
		if !withoutMoving {
			if fc.moving != nil {
				min.SetDelta(0, fc.moving.Y())
			}
		}
		max := lc.Bound().Max
		if !withoutMoving {
			if lc.moving != nil {
				max.SetDelta(0, lc.moving.Y())
			}
		}
		return g.NewBound(&min, &max)
	}
	return nil
}

// SetScrollBarPosition ...
func (o *ScrollViewList) SetScrollBarPosition() {
	// 最新のリストBoundを取得
	lb := o.ScrollBound(false)
	lp, ls := lb.ToPosSize()
	_, ss := o.Bound().ToPosSize()
	sh := ss.H() - (int(LineSpacing) * 2)
	by := math.Abs(lp.Y() * float64(sh) / float64(ls.H()))

	//
	bar := o.Rel("bar").(interfaces.UIControl)
	slider := bar.Rel("slider").(interfaces.UIControl)
	pos, size := slider.Bound().ToPosSize()
	slider.Bound().SetPosSize(g.NewPoint(pos.X(), by+LineSpacing), size)
}

// ScrollViewRow ...
type ScrollViewRow struct {
	UIControl
	index     int
	baseColor *color.RGBA
	fontSet   *char.Resource
	ds        []interface{}
	dsImages  []image.Image
}

// GetControls ...
func (o *ScrollViewRow) GetControls() []interfaces.UIControl {
	if o._childrenCache != nil {
		return o._childrenCache
	}
	ret := _getControls(o)
	o._childrenCache = ret
	return ret
}

// NewScrollViewRow ...
func NewScrollViewRow(s interfaces.Scene, dataSet []interface{}) *ScrollViewRow {
	// 行
	img := utils.CreateRectImage(1, 1, &color.RGBA{64, 64, 64, 127})
	b := &UIControl{
		t:            enum.ControlTypeDefault,
		label:        "scrollview-row",
		bound:        *g.NewBoundByPosSize(g.NewPoint(0, 0), g.NewSize(492, 25)),
		scale:        *g.DefScale(),
		colorScale:   *g.DefCS(),
		scene:        s,
		relations:    map[string]interface{}{},
		eventHandler: event.NewEventHandler(),
		drawer:       NewDefaultDrawer(ebiten.NewImageFromImage(img)),
	}
	o := &ScrollViewRow{
		UIControl: *b,
		baseColor: &color.RGBA{192, 192, 192, 127},
	}

	o.fontSet = char.Res.Get(12, enum.FontStyleGenShinGothicBold)
	o.ds = dataSet

	o.dsImages = make([]image.Image, len(dataSet))
	for i, ds := range dataSet {
		o.dsImages[i] = CreateImageByDataSource(ds, o.fontSet)
	}

	return o
}

// Update ...
func (o *ScrollViewRow) Update() error {
	_ = o.UIControl.Update()

	// 親（ScrollViewList）範囲外を描画対象から外す
	bound := g.NewBound(&o.bound.Min, &o.bound.Max)
	bound.SetDelta(o.moving, nil)
	regionBound := o.parent.Bound()

	if bound.Max.Y() < 0 {
		// 行全体が上に出てる
		o.drawer.withoutDraw = true
	} else if bound.Min.Y() > regionBound.Max.Y()-regionBound.Min.Y() {
		// 行全体が下に出てる
		o.drawer.withoutDraw = true
	} else if bound.Min.Y() < 0 {
		// 行の一部が上に出てる
		_, size := bound.ToPosSize()

		subRect := g.NewBound(g.NewPoint(bound.Min.X(), -bound.Min.Y()), g.NewPoint(float64(size.W()), float64(size.H())))
		o.drawer.subImageRect = subRect.ToImageRect()
		o.drawer.position.SetDelta(0, LineSpacing-bound.Min.Y())
		// log.Printf("行の一部が上に出てる: %#v", subRect)
	} else if bound.Max.Y() > regionBound.Max.Y()-regionBound.Min.Y() {
		// 行の一部が下に出てる
		_, size := bound.ToPosSize()

		maxY := bound.Max.Y() - (regionBound.Max.Y() - regionBound.Min.Y())
		subRect := g.NewBound(g.DefPoint(), g.NewPoint(bound.Max.X(), float64(size.H())-maxY))
		o.drawer.subImageRect = subRect.ToImageRect()
		// log.Printf("行の一部が下に出てる: maxY=%0.2f, bound.max.Y=%0.2f, regionBound.Max.Y()=%0.2f, regionBound.Min.Y()=%0.2f, subRect=%#v", maxY, bound.Max.Y(), regionBound.Max.Y(), regionBound.Min.Y(), subRect)
	}

	return nil
}
