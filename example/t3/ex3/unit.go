package ex3

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/myanagisawa/ebitest/utils"
	"golang.org/x/image/draw"
)

type (
	// Unit ...
	Unit interface {
		Scene
		Collision(u *Unit)
		GetCenter() (float64, float64)
		GetEntity() *Circle
		GetRader() *Circle
		SetCaptured(units []Unit)
	}

	// UnitImpl ...
	UnitImpl struct {
		label     string
		entity    *Circle
		x         float64
		y         float64
		angle     int
		speed     int
		collision Unit
		captured  []Unit
		rader     *Circle
		parent    *Game
	}
)

var (
	maxAngle = 360
)

// NewMyUnit ...
func NewMyUnit(parent *Game) (Unit, error) {
	rand.Seed(time.Now().UnixNano()) //Seed
	// mask画像読み込み
	// mask, _ := utils.GetImageByPath("resources/system_images/mask.png")
	// http://tech.nitoyon.com/ja/blog/2015/12/31/go-image-gen/
	// 座標が円に入っているか
	// http://imagawa.hatenadiary.jp/entry/2016/12/31/190000

	r := 20
	// ユニット画像読み込み
	eimg := getImage("unit-1", r*2, r*2)
	e := &Circle{r: r, image: *eimg}

	unitImpl := &UnitImpl{
		label:  "myUnit",
		entity: e,
		x:      float64(300),
		y:      float64(400),
		angle:  32,
		speed:  4,
		parent: parent,
	}
	// unitImpl.updatePoint()

	r = 100
	// 索敵範囲画像読み込み
	eimg = getImage("search-1", r*2, r*2)

	area := &Circle{r: r, image: *eimg}
	unitImpl.rader = area

	return unitImpl, nil
}

// NewUnit ...
func NewUnit(parent *Game) (Unit, error) {
	rand.Seed(time.Now().UnixNano()) //Seed

	r := 50
	// ユニット画像読み込み
	eimg := getImage("unit-2", r*2, r*2)
	e := &Circle{r: r, image: *eimg}

	unitImpl := &UnitImpl{
		label:  "coin2",
		entity: e,
		x:      float64(600),
		y:      float64(200),
		angle:  270,
		speed:  0,
		parent: parent,
	}
	// unitImpl.updatePoint()

	return unitImpl, nil
}

// NewDebris ...
func NewDebris(speed int, parent *Game) (Unit, error) {
	rand.Seed(time.Now().UnixNano()) //Seed

	rd, gr, bl := uint8(rand.Intn(55)+200), uint8(rand.Intn(55)+200), uint8(rand.Intn(55)+200)

	r := rand.Intn(80) + 20
	// 指定した円の画像を作成
	eimg := createCircleImage(r, color.RGBA{rd, gr, bl, 255}, color.RGBA{0, 0, 0, 255})
	e := &Circle{r: r, image: *eimg}

	x, y := float64(rand.Intn(parent.WindowSize.Width-e.r)), float64(rand.Intn(parent.WindowSize.Height-e.r))
	if int(x) < e.r {
		x = float64(e.r)
	}
	if int(y) < e.r {
		y = float64(e.r)
	}

	// ebitenのrotateとtranslateはy軸0が最上段なので注意
	a := rand.Intn(maxAngle)
	log.Printf("angle: %d, speed: %d", a, speed)

	unitImpl := &UnitImpl{
		entity: e,
		angle:  a,
		x:      x,
		y:      y,
		speed:  speed,
		parent: parent,
	}
	// unitImpl.updatePoint()

	return unitImpl, nil
}

// Update ...
func (s *UnitImpl) Update() error {
	vx, vy := getMoveAmount(s.angle, s.speed)
	s.x += vx
	s.y -= vy

	w := s.parent.WindowSize.Width
	if s.Left() < 0 || w <= s.Right() {
		s.angle = 180 - s.angle
		// s.updatePoint()
	}
	h := s.parent.WindowSize.Height
	if s.Top() < 0 || h <= s.Bottom() {
		s.angle = 360 - s.angle
		// s.updatePoint()
	}

	if s.collision != nil {
		// 衝突時の方向更新
		a := s.angle
		if s.angle > 180 {
			a = s.angle - 360
		}
		s.angle = 180 + a
		// 位置を戻す
		s.x -= vx * 2
		s.y += vy * 2

		s.collision = nil
	}

	return nil
}

// Draw ...
func (s *UnitImpl) Draw(r *ebiten.Image) {

	// 描画オプション: 中心基準に移動、中心座標で回転
	w, h := s.entity.image.Size()
	x, y := s.GetCenter()
	op := defaultDrawOption(x, y, w, h, s.angle)
	if s.collision != nil {
		op.ColorM.Scale(0.5, 0.5, 0.5, 1.0)
	}
	r.DrawImage(&s.entity.image, op)

	if s.captured != nil {
		// draw line
		for _, u := range s.captured {

			x1, y1 := s.GetCenter()
			x2, y2 := u.GetCenter()
			sx, lx, sy, ly := math.Min(x1, x2), math.Max(x1, x2), math.Min(y1, y2), math.Max(y1, y2)
			x, y := lx-sx, ly-sy
			n := math.Atan2(y, x)
			d := n * 180 / math.Pi
			log.Printf("atan2(%f, %f)=%f, deg=%f", y, x, n, d)
			log.Printf("  sx=%f, lx=%f, sy=%f, ly=%f, ", sx, lx, sy, ly)

			ebitenutil.DrawLine(r, x1, y1, x2, y2, color.RGBA{0, 255, 0, 255})
		}
		s.captured = nil
	}

	// 索敵範囲を描画
	if s.rader != nil {
		drawRader(s, r)
	}
	drawUnitCircle(s, r)
	// drawHikakuCircle(s, r)
}

// defaultDrawOption デフォルト描画オプション
func defaultDrawOption(x, y float64, w, h, a int) *ebiten.DrawImageOptions {
	// 描画オプション: 中心基準に移動、中心座標で回転
	op := &ebiten.DrawImageOptions{}
	// 描画位置指定
	op.GeoM.Reset()
	// 対象画像の縦横半分だけマイナス位置に移動（原点に中心座標が来るように移動する）
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	// 中心を軸に回転
	op.GeoM.Rotate(-2 * math.Pi * float64(a) / float64(maxAngle))
	// ユニットの座標に移動
	op.GeoM.Translate(x, y)

	return op
}

func drawRader(s *UnitImpl, r *ebiten.Image) {
	// 描画オプション: 中心基準に移動、中心座標で回転
	w, h := s.rader.image.Size()
	x, y := s.GetCenter()
	op := defaultDrawOption(x, y, w, h, s.angle)

	op.ColorM.Scale(1.0, 1.0, 1.0, 0.25)
	r.DrawImage(&s.rader.image, op)
}

func drawUnitCircle(s *UnitImpl, r *ebiten.Image) {
	// 指定した円の画像を作成
	eimg := createCircleImage(s.entity.r, color.RGBA{4, 124, 208, 128}, color.RGBA{143, 215, 212, 128})

	// 描画オプション: 中心基準に移動、中心座標で回転
	w, h := eimg.Size()
	x, y := s.GetCenter()
	op := defaultDrawOption(x, y, w, h, s.angle)

	r.DrawImage(eimg, op)
}

// GetSize unitのサイズを返します
func (s *UnitImpl) GetSize() (int, int) {
	return s.entity.r * 2, s.entity.r * 2
}

// GetCenter unitの中心座標を返します
func (s *UnitImpl) GetCenter() (float64, float64) {
	return s.x, s.y
}

// GetEntity unitのentityを返します
func (s *UnitImpl) GetEntity() *Circle {
	return s.entity
}

// GetRader ...
func (s *UnitImpl) GetRader() *Circle {
	return s.rader
}

// Collision ...
func (s *UnitImpl) Collision(c *Unit) {
	s.collision = *c
}

// SetCaptured ...
func (s *UnitImpl) SetCaptured(units []Unit) {
	s.captured = units
}

// distance x, yが表す点が半径rの円の範囲内に位置する場合、1以下、範囲外の場合1以上を返します
func distance(x, y, r int) float64 {
	var dx, dy int = int(r) - x, int(r) - y
	return math.Sqrt(float64(dx*dx+dy*dy)) / float64(r)
}

// Left ...
func (s *UnitImpl) Left() int {
	return int(s.x) - s.entity.r
}

// Top ...
func (s *UnitImpl) Top() int {
	return int(s.y) - s.entity.r
}

// Right ...
func (s *UnitImpl) Right() int {
	return int(s.x) + s.entity.r
}

// Bottom ...
func (s *UnitImpl) Bottom() int {
	return int(s.y) + s.entity.r
}

// Width ...
func (s *UnitImpl) Width() int {
	return 2 * s.entity.r
}

// Height ...
func (s *UnitImpl) Height() int {
	return 2 * s.entity.r
}

func getMoveAmount(angle, speed int) (vx, vy float64) {
	rad := float64(angle) * math.Pi / 180
	vx = math.Cos(rad) * float64(speed)
	vy = math.Sin(rad) * float64(speed)
	return vx, vy
}

// func (s *UnitImpl) updatePoint() {
// 	rad := float64(s.angle) * math.Pi / 180
// 	s.vx = math.Cos(rad) * float64(s.speed)
// 	s.vy = math.Sin(rad) * float64(s.speed)
// 	// if s.label == "myCoin" {
// 	// 	log.Printf("rad=%f, vx=%f, vy=%f, sin(rad)=%f", rad, s.vx, s.vy, math.Sin(rad))
// 	// }
// }

// createCircleImage 半径rの円の画像イメージを作成します。color1は円の色、color2は円の向きを表す線の色です
func createCircleImage(r int, color1, color2 color.RGBA) *ebiten.Image {
	m := image.NewRGBA(image.Rect(0, 0, r*2, r*2))
	// 横ループ、半径*2＝直径まで
	for x := 0; x < r*2; x++ {
		// 縦ループ、半径*2＝直径まで
		for y := 0; y < r*2; y++ {
			// 向き先判定中心より右側の中心から水平2pixel分
			if x > r && y >= r-1 && y <= r+1 {
				m.Set(x, y, color2)
			} else {
				d := distance(x, y, r)
				if d > 1 {
					// 円の範囲外の点は透明で描画
					m.Set(x, y, color.RGBA{0, 0, 0, 0})
				} else {
					// 円の範囲内の点を指定された色で描画
					m.Set(x, y, color1)
				}
			}
		}
	}
	eimg, err := ebiten.NewImageFromImage(m, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	return eimg
}

// getImage 指定した名称の画像を読み込みます(w, h: 縦横サイズ)
func getImage(name string, w, h int) *ebiten.Image {
	path := fmt.Sprintf("resources/system_images/%s.png", name)
	img, err := utils.OrientationImage(path)
	if err != nil {
		panic(err)
	}
	log.Printf("img.Bounds: %#v", img.Bounds())

	// リサイズ
	imgDst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.CatmullRom.Scale(imgDst, imgDst.Bounds(), img, img.Bounds(), draw.Over, nil)

	// img, err := utils.ScaleImage(img, w, h)
	// if err != nil {
	// 	panic(err)
	// }
	eimg, err := ebiten.NewImageFromImage(imgDst, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	return eimg
}
