package ex3

const ()

type (
	// Window ...
	Window struct {
		Rows []Scene
	}

	// UnitInfo ...
	UnitInfo struct {
		Unit *Unit
	}
)

var ()

// Update ...
func (s *UnitInfo) Update() error {
	return nil
}

// // Draw ...
// func (s *UnitInfo) Draw(r *ebiten.Image) {
// 	ebitenutil.DrawRect(r, x-float64(s.entity.r), y+float64(s.entity.r), w, 5, color.RGBA{0, 255, 0, 255})

// }
