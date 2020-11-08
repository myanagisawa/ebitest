package scene

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/myanagisawa/ebitest/example/t5/ebitest"
	"github.com/myanagisawa/ebitest/example/t5/enum"
	"github.com/myanagisawa/ebitest/example/t5/interfaces"
	"github.com/myanagisawa/ebitest/example/t5/models/control"
	"github.com/myanagisawa/ebitest/example/t5/models/layer"
)

type (
	// MainMenu ...
	MainMenu struct {
		Base
	}
)

// NewMainMenu ...
func NewMainMenu(m interfaces.GameManager) *MainMenu {

	s := &MainMenu{
		Base: Base{
			label: "MainMenu",
		},
	}

	l := layer.NewLayerBase("Layer1", ebitest.Images["bgFlower"], s, ebitest.NewScale(1.0, 1.0), nil, 0, false)
	s.SetLayer(l)

	c := control.NewButton("マップ", l, ebitest.Images["btnBase"], ebitest.Fonts["btnFont"], color.Black, 500, 500)
	l.AddUIControl(c)
	l.EventHandler().AddEventListener(c, "click", func(target interfaces.UIControl, scene interfaces.Scene, point *ebitest.Point) {
		log.Printf("%s clicked", target.Label())
		m.TransitionTo(enum.MapEnum)
	})

	img := ebitest.CreateRectImage(400, 600, color.RGBA{0, 0, 0, 128})
	l = layer.NewLayerBase("Layer2", img, s, ebitest.NewScale(0.7, 0.7), ebitest.NewPoint(10.0, 50.0), 0, true)
	s.SetLayer(l)

	c = control.NewButton("Battle(dev)", l, ebitest.Images["btnBase"], ebitest.Fonts["btnFont"], color.Black, 50, 100)
	l.AddUIControl(c)
	l.EventHandler().AddEventListener(c, "click", func(target interfaces.UIControl, scene interfaces.Scene, point *ebitest.Point) {
		log.Printf("%s clicked", target.Label())
		layer := scene.GetLayerByLabel("Layer1")
		s := layer.EbiObjects()[0].Scale()
		if s.X() >= 1.0 {
			s.Set(0.5, 0.5)
		} else {
			s.Set(1.0, 1.0)
		}

	})
	// l.AddEventListener(c, "click", func(target UIControl, source *EventSource) {
	// 	log.Printf("btnBattle clicked")
	// 	source.scene.Manager().TransitionToBattleScene()
	// })

	img = ebitest.CreateRectImage(600, 400, color.RGBA{255, 32, 32, 128})
	l = layer.NewLayerBase("Layer3", img, s, nil, ebitest.NewPoint(200.0, 100.0), 0, false)
	s.SetLayer(l)

	c = control.NewButton("scale", l, ebitest.Images["btnBase"], ebitest.Fonts["btnFont"], color.Black, 50, 30)
	l.AddUIControl(c)
	l.EventHandler().AddEventListener(c, "click", func(target interfaces.UIControl, scene interfaces.Scene, point *ebitest.Point) {
		log.Printf("%s clicked", target.Label())
		layer := scene.GetLayerByLabel("Layer3")
		s := layer.EbiObjects()[0].Scale()
		if s.X() >= 1.0 {
			s.Set(0.5, 0.5)
		} else {
			s.Set(1.0, 1.0)
		}
	})

	c = control.NewText("【公式】タイムマシーン3号 漫才 「お見合い」Ayi|`", l, ebitest.ScaleFonts[8], color.Black, 10, 80)
	l.AddUIControl(c)

	c = control.NewText("【公式】タイムマシーン3号 漫才 「お見合い」Ayi|`", l, ebitest.ScaleFonts[16], color.Black, 10, 120)
	l.AddUIControl(c)

	c = control.NewColumn("【公式】タイムマシーン3号 漫才 「お見合い」Ayi|`", l, ebitest.ScaleFonts[24], color.Black, color.RGBA{200, 200, 200, 128}, 10, 160)
	l.AddUIControl(c)

	img = ebitest.CreateRectImage(800, 300, color.RGBA{128, 128, 128, 192})
	l = layer.NewLayerBase("Layer4", img, s, nil, ebitest.NewPoint(300.0, 550.0), 0, false)
	s.SetLayer(l)

	// スクロールビュー
	cols := []interface{}{
		"ID", "カラム名ですColumn", "3番目", "HOGE",
	}
	data := [][]interface{}{
		{1, "１行目", "あいうえお", 12345},
		{2, "２行目", "かきくけこ", 12345},
		{3, "３行目", "さしすせそ", 12345},
		{4, "４行目", "たちつてと", 12345},
		{5, "５行目", "なにぬねの", 12345},
		{6, "６行目", "はひふへほ", 12345},
		{7, "７行目", "まみむめも", 12345},
		{8, "８行目", "やいゆえよ", 12345},
		{9, "９行目", "らりるれろ", 12345},
		{10, "１０行目", "わをん", 12345},
		{11, "１１行目", "あいうえお", 12345},
		{12, "１２行目", "かきくけこ", 12345},
		{13, "１３行目", "さしすせそ", 12345},
		{14, "１４行目", "たちつてと", 12345},
		{15, "１５行目", "なにぬねの", 12345},
		{16, "１６行目", "はひふへほ", 12345},
		{17, "１７行目", "まみむめも", 12345},
		{18, "１８行目", "やいゆえよ", 12345},
		{19, "１９行目", "らりるれろ", 12345},
		{20, "２０行目", "わをん", 12345},
	}
	c = control.NewUIScrollViewByList(l, cols, data, 700, 250, 30, ebitest.NewPoint(20.0, 20.0))
	l.AddUIControl(c)
	l.EventHandler().AddEventListener(c, "click", func(target interfaces.UIControl, scene interfaces.Scene, point *ebitest.Point) {
		// log.Printf("%s clicked", target.Label())
		t := target.(*control.UIScrollViewImpl)
		index := t.GetIndexOfFocusedRow()
		log.Printf("clicked: index: %d, data: %#v", index, data[index])
	})

	return s
}

// Update ...
func (s *MainMenu) Update() error {
	log.Printf("(s *MainMenu) Update()")
	et := GetEdgeType(ebiten.CursorPosition())
	if et != enum.EdgeTypeNotEdge {
		s.layers[0].Scroll(et)
	}

	s.activeLayer = s.LayerAt(ebiten.CursorPosition())
	if s.activeLayer != nil {
		// log.Printf("activeLayer: %#v", s.activeLayer.Label())
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			// click イベントを発火
			s.activeLayer.EventHandler().Firing(s, "click", x, y)
		}
	}

	for _, layer := range s.layers {
		layer.Update()
	}

	return nil
}

// Draw ...
func (s *MainMenu) Draw(screen *ebiten.Image) {
	log.Printf("(s *MainMenu) Draw(screen *ebiten.Image)")

	for _, layer := range s.layers {
		layer.Draw(screen)
	}

	active := " - "
	control := " - "
	if s.activeLayer != nil {
		eo := s.activeLayer.EbiObjects()[0]
		px, py := eo.GlobalPosition().Get()
		active = fmt.Sprintf("%s: (%d, %d)", s.activeLayer.LabelFull(), int(px), int(py))
		c := s.activeLayer.UIControlAt(ebiten.CursorPosition())
		if c != nil {
			px, py = c.EbiObjects()[0].GlobalPosition().Get()
			control = fmt.Sprintf("%s: (%d, %d)", c.Label(), int(px), int(py))
		}
	}

	// mask := ebiten.NewImage(200, 200)
	// mask.Fill(color.White)

	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(50, 650)
	// screen.DrawImage(mask.SubImage(image.Rect(0, 0, 100, 100)).(*ebiten.Image), op)

	// i1 := ebiten.NewImage(100, 100)
	// i1.Fill(color.Black)

	// op = &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(50, 50)
	// mask.DrawImage(i1, op)

	x, y := ebiten.CursorPosition()
	dbg := fmt.Sprintf("%s\nTPS: %0.2f\nFPS: %0.2f\npos: (%d, %d)\nactive:\n - layer: %s\n - control: %s", printMemoryStats(), ebiten.CurrentTPS(), ebiten.CurrentFPS(), x, y, active, control)
	ebitenutil.DebugPrint(screen, dbg)
}
