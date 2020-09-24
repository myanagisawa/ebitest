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
	"github.com/hajimehoshi/ebiten/text"
)

type (
	// Unit ...
	Unit interface {
		Scene
		Belongs() int
		UpdateHP(damage int)
		Collision(u *Unit)
		GetCenter() (float64, float64)
		GetEntity() *Circle
		GetRader() *Circle
		SetCaptured(units []Unit)
		GetStatus() int
		SetFocus(b bool)
	}

	// UnitImpl ...
	UnitImpl struct {
		Label     string
		HP        int
		MaxHP     int
		belongs   int
		entity    *Circle
		x         float64
		y         float64
		angle     int
		speed     int
		collision Unit
		captured  []Unit
		locked    Unit
		rader     *Circle
		parent    *Scene
		status    int
		infoList  []*DamageLabel
		focused   bool
	}

	// DamageLabel ...
	DamageLabel struct {
		label  string
		face   *LabelFace
		count  int
		erased int
	}
)

var (
	maxAngle            = 360
	capturedLineCounter = 0
)

// NewDamageLabel ...
func NewDamageLabel(d int, e int) *DamageLabel {
	var ff *LabelFace
	if d > 0 {
		ff = fface10Red
	} else {
		ff = fface10White
	}

	return &DamageLabel{
		label:  fmt.Sprintf("%d", d),
		face:   ff,
		count:  0,
		erased: e,
	}
}

// NewUnit ...
func NewUnit(parent Scene, team, hp, r int, label string, x, y, angle, speed, rader int) (Unit, error) {
	rand.Seed(time.Now().UnixNano()) //Seed

	// ユニット画像読み込み
	img := images[fmt.Sprintf("unit-%d", team)]
	eimg, _ := ebiten.NewImageFromImage(ResizeImage(img, r*2, r*2), ebiten.FilterDefault)
	e := &Circle{r: r, image: eimg}

	unitImpl := &UnitImpl{
		Label:  label,
		HP:     hp,
		MaxHP:  hp,
		entity: e,
		x:      float64(x),
		y:      float64(y),
		angle:  angle,
		speed:  speed,
		parent: &parent,
	}
	// 索敵範囲画像読み込み
	// eimg = getImage("search-1.png", rader*2, rader*2)
	img = images["search-1"]
	eimg, _ = ebiten.NewImageFromImage(ResizeImage(img, rader*2, rader*2), ebiten.FilterDefault)
	area := &Circle{r: rader, image: eimg}
	unitImpl.rader = area

	return unitImpl, nil
}

// Update ...
func (s *UnitImpl) Update() error {
	// 非生存
	if s.status != 0 {
		return nil
	}
	vx, vy := getMoveAmount(s.angle, s.speed)
	s.x += vx
	s.y -= vy

	w, h := (*s.parent).GetSize()
	if s.Left() < 0 || w <= s.Right() {
		s.angle = 180 - s.angle
		// s.updatePoint()
	}
	if s.Top() < 0 || h <= s.Bottom() {
		s.angle = 360 - s.angle
		// s.updatePoint()
	}

	if s.collision != nil {
		// 衝突時の方向更新
		// a := s.angle
		// if s.angle > 180 {
		// 	a = s.angle - 360
		// }
		// s.angle = 180 + a
		s.collision.UpdateHP(1)
		// 位置を戻す
		s.x -= vx * 10
		s.y += vy * 10

		s.collision = nil
	}

	// ロックオン対象を更新
	s.locked = nil
	dist := 0.0
	for _, u := range s.captured {
		// 自機との距離を算出
		x1, y1 := s.GetCenter()
		x2, y2 := u.GetCenter()
		dx, dy := x2-x1, y2-y1
		d := math.Sqrt(float64(dx*dx + dy*dy))
		if s.locked == nil || dist > d {
			s.locked = u
			dist = d
		}
	}
	if capturedLineCounter == 10 {
		capturedLineCounter = -10
	} else {
		capturedLineCounter++
	}
	// log.Printf("locled=%d", s.locked.GetEntity().r)

	if s.locked != nil {
		// ロックオン対象が存在する場合はそちらに方向転換する

		//自機と対象の角度を算出
		x1, y1 := s.GetCenter()
		x2, y2 := s.locked.GetCenter()
		dx, dy := x2-x1, -(y2 - y1) // 画面の上側をY座標＋とするので、Y座標は符号を入れ替える
		// radianを取得
		n := math.Atan2(dy, dx)
		// radian ->degreeに変換
		d := n * 180 / math.Pi
		if s.angle != int(d) {
			// 自機の向きを更新
			// log.Printf("[%s] angle=%d, d=%d, rad=%f", s.label, s.angle, int(d), n)
			if s.angle > int(d) {
				if s.angle-5 < int(d) {
					s.angle = int(d)
				} else {
					s.angle -= 7
				}
			} else if s.angle < int(d) {
				if s.angle+5 < int(d) {
					s.angle = int(d)
				} else {
					s.angle += 7
				}
			}
		}
	}

	// log.Printf("updated unit: %s", s.label)
	return nil
}

// Draw ...
func (s *UnitImpl) Draw(r *ebiten.Image) {
	c := color.RGBA{0, 255, 0, 255}
	if s.belongs == 0 {
		c = color.RGBA{0, 0, 255, 127}
	} else if s.belongs == 1 {
		c = color.RGBA{255, 0, 0, 127}
	} else if s.belongs == 2 {
		c = color.RGBA{0, 255, 0, 127}
	}

	// 描画オプション: 中心基準に移動、中心座標で回転
	w, h := s.entity.image.Size()
	x, y := s.GetCenter()
	// log.Printf("%0.2f", scale)
	op := defaultDrawOption(x, y, float64(w), float64(h), float64(s.angle))
	if s.status == -1 {
		// 行動不能
		op.ColorM.Scale(0.5, 0.5, 0.5, 0.5)
	}
	r.DrawImage(s.entity.image, op)

	if s.captured != nil {
		// draw line
		for _, u := range s.captured {

			x1, y1 := s.GetCenter()
			x2, y2 := u.GetCenter()

			if capturedLineCounter >= 0 {
				ebitenutil.DrawLine(r, x1, y1, x2, y2, c)
			}
			// debug 2点間の距離を表示
			// sqrt( (x1-x2)^2 + (y1-y2)^2 )
			// dx, dy := x1-x2, y1-y2
			// dist := math.Sqrt(float64(dx*dx + dy*dy))
			// log.Printf("[%s] x1=%f, y1=%f, x2=%f, y2=%f: distance=%f", s.label, x1, y1, x2, y2, dist)
			// debug 2点間の角度を表示
			// sx, lx, sy, ly := math.Min(x1, x2), math.Max(x1, x2), math.Min(y1, y2), math.Max(y1, y2)
			// x, y := lx-sx, ly-sy
			// // radianを取得
			// n := math.Atan2(y, x)
			// // radian ->degreeに変換
			// d := n * 180 / math.Pi
			// log.Printf("[%s] degree=%f", s.label, d)
		}
		s.captured = nil
	}

	// 索敵範囲を描画
	if s.rader != nil {
		drawRader(s, r)
	}

	rs := int(float64(s.HP) / float64(s.MaxHP) * 100)
	// log.Printf("hp: %d / %d, rs=%d", s.hp, s.maxHp, rs)
	if rs == 100 {
		ebitenutil.DrawRect(r, x-float64(s.entity.r), y+float64(s.entity.r), float64(s.entity.r)*2, 5, color.RGBA{0, 255, 0, 255})
	} else {
		w := float64(s.entity.r) * 2 * (float64(rs) / 100)
		ebitenutil.DrawRect(r, x-float64(s.entity.r), y+float64(s.entity.r), w, 5, color.RGBA{0, 255, 0, 255})
		ebitenutil.DrawRect(r, x-float64(s.entity.r)+w, y+float64(s.entity.r), (float64(s.entity.r)*2)-w, 5, color.RGBA{127, 127, 127, 127})
	}
	ff := fface10White
	if s.focused {
		ff = fface10Green
	}
	text.Draw(r, s.Label, ff.uiFont, int(x)-10, int(y)-20, ff.uiFontColor)

	// ダメージ表示を描画
	for _, info := range s.infoList {
		info.count++
		text.Draw(r, info.label, info.face.uiFont, int(x), int(y)-(info.count*2), info.face.uiFontColor)
	}
	// infoListを更新
	var l []*DamageLabel
	for _, d := range s.infoList {
		if d.count < d.erased {
			l = append(l, d)
		}
	}
	s.infoList = l

	// drawUnitCircle(s, r)
}

// defaultDrawOption デフォルト描画オプション
func defaultDrawOption(x, y, w, h, a float64) *ebiten.DrawImageOptions {
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
	op := defaultDrawOption(x, y, float64(w), float64(h), float64(s.angle))

	op.ColorM.Scale(1.0, 1.0, 1.0, 0.1)
	r.DrawImage(s.rader.image, op)
}

// func drawUnitCircle(s *UnitImpl, r *ebiten.Image) {
// 	// 指定した円の画像を作成
// 	eimg := createCircleImage(s.entity.r, color.RGBA{4, 124, 208, 128}, color.RGBA{143, 215, 212, 128})

// 	// 描画オプション: 中心基準に移動、中心座標で回転
// 	w, h := eimg.Size()
// 	x, y := s.GetCenter()
// 	op := defaultDrawOption(x, y, w, h, s.angle)

// 	r.DrawImage(eimg, op)
// }

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

// Belongs ...
func (s *UnitImpl) Belongs() int {
	return s.belongs
}

// UpdateHP ...
func (s *UnitImpl) UpdateHP(damage int) {
	// log.Printf("%s: damage: %d", s.label, damage)
	s.infoList = append(s.infoList, NewDamageLabel(damage, 20))
	if s.HP <= damage {
		s.HP = 0
		s.dead()
	} else {
		s.HP -= damage
	}
}

// GetStatus ...
func (s *UnitImpl) GetStatus() int {
	return s.status
}

// SetFocus ...
func (s *UnitImpl) SetFocus(b bool) {
	s.focused = b
}

// Collision ...
func (s *UnitImpl) Collision(c *Unit) {
	s.collision = *c
}

// SetCaptured ...
func (s *UnitImpl) SetCaptured(units []Unit) {
	if units == nil {
		s.captured = nil
		return
	}
	s.captured = append(s.captured, units...)
}

// distance x, yが表す点が半径rの円の範囲内に位置する場合、1以下、範囲外の場合1以上を返します
func distance(x, y, r int) float64 {
	var dx, dy int = r - x, r - y
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

// dead ...
func (s *UnitImpl) dead() {
	s.rader = nil
	s.captured = nil
	s.collision = nil
	s.status = -1
	log.Printf("%s: 行動不能", s.Label)
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
