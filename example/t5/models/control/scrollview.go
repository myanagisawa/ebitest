package control

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/myanagisawa/ebitest/example/t5/ebitest"
	"github.com/myanagisawa/ebitest/example/t5/interfaces"
	"github.com/myanagisawa/ebitest/example/t5/models"
	"github.com/myanagisawa/ebitest/example/t5/models/input"
	"github.com/myanagisawa/ebitest/utils"
)

const (
	margin  = 2
	padleft = 2
)

type listRowView struct {
	UIControlImpl
	source   []interface{}
	texts    []*ebiten.Image
	cols     []*ebiten.Image
	colWidth []float64
}

// newListRowView ...
func newListRowView(row []interface{}, w, h int) *listRowView {
	eimg, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)

	label := fmt.Sprintf("row-%s", utils.RandomLC(8))
	con := &UIControlImpl{
		label:          label,
		bg:             models.NewEbiObject(label, eimg, nil, nil, ebitest.NewPoint(0, 0), 0, true, true, false),
		hasHoverAction: true,
	}

	texts := []*ebiten.Image{}
	var timg *ebiten.Image
	for i := range row {
		col := row[i]
		switch val := col.(type) {
		case string:
			t := utils.CreateTextImage(val, ebitest.Fonts["btnFont"], color.Black)
			timg, _ = ebiten.NewImageFromImage(*t, ebiten.FilterDefault)
			texts = append(texts, timg)
			// log.Printf("string: %s", val)
		case int:
			t := utils.CreateTextImage(fmt.Sprintf("%d", val), ebitest.Fonts["btnFont"], color.Black)
			timg, _ = ebiten.NewImageFromImage(*t, ebiten.FilterDefault)
			texts = append(texts, timg)
			// log.Printf("int: %d", val)
		}
	}

	return &listRowView{
		UIControlImpl: *con,
		source:        row,
		texts:         texts,
		cols:          []*ebiten.Image{},
	}
}

func (r *listRowView) GetRowHeight() int {
	_, h := r.bg.Size()
	return h
}

// DrawListRow ...
func (r *listRowView) DrawListRow() *ebiten.Image {
	eimg := r.bg.EbitenImage()
	_, bh := eimg.Size()

	x := 0.0
	for i, row := range r.texts {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, margin)
		eimg.DrawImage(r.cols[i], op)

		_, h := row.Size()

		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x+padleft, float64(bh-h)/2)
		eimg.DrawImage(row, op)

		x += r.colWidth[i] + margin
	}

	return eimg
}

type listView struct {
	UIControlImpl
	rows []listRowView
}

// newListView ...
func newListView(data [][]interface{}, roww, rowh int) *listView {
	listh := 0
	rows := []listRowView{}
	for i := range data {
		rowview := newListRowView(data[i], roww, rowh)
		rows = append(rows, *rowview)
		listh += rowh
	}
	// margin分を高さに追加
	listh += margin * (len(data) - 1)

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
			img := ebitest.CreateBorderedRectImage(int(colWidth[i]), rowh-margin, color.RGBA{255, 32, 32, 64})
			eimg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
			colBg[i] = eimg
		}

		// rowViewにカラムサイズリストをセット
		for i := range rows {
			row := rows[i]
			row.colWidth = colWidth
			row.cols = colBg
			rows[i] = row
		}

	}

	// ベース画像を作成
	eimg, _ := ebiten.NewImage(roww, listh, ebiten.FilterDefault)

	label := fmt.Sprintf("list-%s", utils.RandomLC(8))
	con := &UIControlImpl{
		label:          label,
		bg:             models.NewEbiObject(label, eimg, nil, nil, ebitest.NewPoint(0, 0), 0, true, true, false),
		hasHoverAction: true,
	}

	l := &listView{
		UIControlImpl: *con,
		rows:          rows,
	}
	return l
}

// DrawList ...
func (l *listView) DrawList() *ebiten.Image {
	eimg := l.bg.EbitenImage()

	y := 0.0
	for _, row := range l.rows {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, y)

		y += float64(row.GetRowHeight()) + margin
		eimg.DrawImage(row.DrawListRow(), op)
	}

	return eimg
}

func (l *listView) GetListHeight() int {
	listHeight := 0
	for i := range l.rows {
		row := l.rows[i]
		listHeight += row.GetRowHeight()
	}
	return listHeight
}

// UIScrollViewImpl ...
type UIScrollViewImpl struct {
	UIControlImpl
	list     *listView
	listRect image.Rectangle
	stroke   *input.Stroke
}

// NewUIScrollViewByList ...
func NewUIScrollViewByList(parent interfaces.Layer, data [][]interface{}, w, h, rowh int) interfaces.UIScrollView {
	// img := ebitest.CreateBorderedRectImage(w, h, color.RGBA{0, 0, 0, 255})
	eimg, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)

	label := fmt.Sprintf("scroll-%s", utils.RandomLC(8))
	con := &UIControlImpl{
		label:          label,
		bg:             models.NewEbiObject(label, eimg, parent.EbiObjects()[0], nil, ebitest.NewPoint(20, 20), 0, true, true, false),
		hasHoverAction: true,
	}

	l := newListView(data, w, rowh)
	s := &UIScrollViewImpl{
		UIControlImpl: *con,
		list:          l,
		listRect:      image.Rect(0, 0, w, h),
	}
	return s
}

// Update ...
func (c *UIScrollViewImpl) Update(screen *ebiten.Image) error {
	// ホイールイベント
	_, dy := ebiten.Wheel()
	// log.Printf("%0.1f < %0.1f && %0.1f > %0.1f", c.contentOffset.y, c.scrollMin, c.contentOffset.y, c.scrollMax)
	if dy < 0 {

		if c.listRect.Min.Y > 0 {
			c.listRect = image.Rect(c.listRect.Min.X, c.listRect.Min.Y+int(dy*2), c.listRect.Max.X, c.listRect.Max.Y+int(dy*2))
			// c.listRect.Add(image.Point{X: 0, Y: int(dy * 5)})
			// log.Printf("listRect: %#v", c.listRect)
		}

	} else {

		// 表示領域高さ
		h := float64(c.listRect.Max.Y - c.listRect.Min.Y)
		// コンテンツサイズ
		_, ch := c.list.bg.Size()
		// スクロール最大位置
		scrollMax := int((float64(ch) * c.list.bg.Scale().Y()) - h)
		if c.listRect.Min.Y < scrollMax {
			c.listRect = image.Rect(c.listRect.Min.X, c.listRect.Min.Y+int(dy*5), c.listRect.Max.X, c.listRect.Max.Y+int(dy*5))
			// c.listRect.Add(image.Point{X: 0, Y: int(dy * 5)})
			// log.Printf("listRect: %#v", c.listRect)
		}

	}

	// c.listRect = image.Rect(c.listRect.Min.X, c.listRect.Min.Y+int(dy*5), c.listRect.Max.X, c.listRect.Max.Y+int(dy*5))
	return nil
}

// Draw ...
func (c *UIScrollViewImpl) Draw(screen *ebiten.Image) {
	list := c.list.DrawList()
	// x, y := list.Size()
	// log.Printf("listsize: %d, %d", x, y)

	// リストの表示部分を描画
	op := &ebiten.DrawImageOptions{}
	w, h := c.bg.Size()
	// 描画位置指定
	op.GeoM.Reset()

	op.GeoM.Scale(c.bg.GlobalScale())

	// 対象画像の縦横半分だけマイナス位置に移動（原点に中心座標が来るように移動する）
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	// 中心を軸に回転
	op.GeoM.Rotate(c.bg.Theta())
	// ユニットの座標に移動
	op.GeoM.Translate(float64(w)/2, float64(h)/2)

	op.GeoM.Translate(c.bg.GlobalPosition())
	screen.DrawImage(list.SubImage(c.listRect).(*ebiten.Image), op)
}
