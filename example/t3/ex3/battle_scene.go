package ex3

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type (
	// BattleScene ...
	BattleScene struct {
		bg      Scene
		size    Size
		teams   []*Team
		windows []Scene
	}

	// Team ...
	Team struct {
		No        int
		Units     []Unit
		Enemies   []*Team
		Alliances []*Team
		Location  *Point
		Parent    Scene
		IsAllies  bool
	}
)

var (
	pins []*Pin
)

// NewBattleScene ...
func NewBattleScene(s Size) *BattleScene {

	backGround, _ := NewBackGround(s)
	bs := &BattleScene{
		bg:   backGround,
		size: s,
	}

	return bs
}

// AddTeam ...
func (s *BattleScene) AddTeam(loc *Point, isAllies bool) *Team {
	t := &Team{No: len(s.teams), Location: loc, Parent: s, IsAllies: isAllies}
	s.teams = append(s.teams, t)
	return s.teams[len(s.teams)-1]
}

// InitWindows ...
func (s *BattleScene) InitWindows() {
	var units []Unit
	for _, team := range s.teams {
		units = append(units, team.Units...)
	}
	is1 := NewInfoScene(units)
	s.windows = append(s.windows, is1)

}

// Update ...
func (s *BattleScene) Update() error {

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		log.Printf("Clicked: %d, %d", x, y)

		img := images["pin"]
		eimg := ebiten.NewImageFromImage(ResizeImage(img, 9, 25))
		pin := &Pin{&Point{x, y}, eimg}
		pins = append(pins, pin)
	}

	if err := s.bg.Update(); err != nil {
		return err
	}

	// 各unitのレーダーと衝突判定処理
	for _, team := range s.teams {
		for _, u := range team.Units {
			if u.GetStatus() != 0 {
				continue
			}
			// ユニットのレーダー捕捉判定
			u.SetCaptured(nil)
			for _, et := range team.Enemies {
				// log.Printf("et: %d, %v", et.No, et.Units)
				captureUnit(u, et.Units)
			}

			// ユニットの衝突判定
			for _, et := range team.Enemies {
				for _, eu := range et.Units {
					if eu.GetStatus() != 0 {
						continue
					}
					// log.Printf("Collision")
					if CollisionUnit(u, eu) {
						u.Collision(&eu)
					}
				}
			}
		}
	}

	// 各unitのupdate処理
	for _, team := range s.teams {
		for _, u := range team.Units {
			if u.GetStatus() != 0 {
				continue
			}
			if err := u.Update(); err != nil {
				return err
			}
		}
	}

	// 各windowのupdate処理
	for _, w := range s.windows {
		if err := w.Update(); err != nil {
			return err
		}
	}

	return nil
}

// Draw ...
func (s *BattleScene) Draw(r *ebiten.Image) {
	s.bg.Draw(r)

	for _, team := range s.teams {
		for _, u := range team.Units {
			u.Draw(r)
		}
	}

	for _, pin := range pins {
		pin.Draw(r)
	}

	for _, w := range s.windows {
		w.Draw(r)
	}
}

// Draw ...
func (p *Pin) Draw(r *ebiten.Image) {
	w, h := p.image.Size()
	// 描画オプション: 中心基準に移動、中心座標で回転
	op := &ebiten.DrawImageOptions{}
	// 描画位置指定
	op.GeoM.Reset()
	// 対象画像の縦横半分だけマイナス位置に移動（原点に中心座標が来るように移動する）
	op.GeoM.Translate(-float64(w), -float64(h))
	// ユニットの座標に移動
	op.GeoM.Translate(float64(p.point.X), float64(p.point.Y))

	r.DrawImage(p.image, op)
}

// GetSize ...
func (s *BattleScene) GetSize() (int, int) {
	return s.size.Width, s.size.Height
}

// CollisionUnit unit同士の衝突状態を返す
func CollisionUnit(unit1, unit2 Unit) bool {
	x1, y1 := unit1.GetCenter()
	x2, y2 := unit2.GetCenter()
	e1, e2 := unit1.GetEntity(), unit2.GetEntity()
	// (xc1-xc2)^2 + (yc1-yc2)^2 ≦ (r1+r2)^2
	var dx, dy, dr float64 = float64(x1 - x2), float64(y1 - y2), float64(e1.R() + e2.R())
	if (dx*dx + dy*dy) <= dr*dr {
		return true
	}
	return false
}

// captureUnit unitの索敵範囲に入ったunitsを取得する
func captureUnit(unit Unit, units []Unit) {
	x1, y1 := unit.GetCenter()
	c1 := unit.GetRader()
	captured := []Unit{}
	for _, u := range units {
		if u.GetStatus() != 0 {
			continue
		}
		x2, y2 := u.GetCenter()
		e := u.GetEntity()

		var dx, dy, dr float64 = float64(x1 - x2), float64(y1 - y2), float64(c1.R() + e.R())
		if (dx*dx + dy*dy) <= dr*dr {
			// レーダー捕捉
			captured = append(captured, u)
			// log.Printf("captured!!")
		}
	}
	unit.SetCaptured(captured)
}
