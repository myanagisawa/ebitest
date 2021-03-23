package wmap

import (
	"log"

	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/lib/game"
	"github.com/myanagisawa/ebitest/example/t7/lib/interfaces"
	"github.com/myanagisawa/ebitest/example/t7/lib/models/control"
)

type (
	// Scene ...
	Scene struct {
		children []interfaces.UIControl
		manager  *game.Manager
	}
)

var (
	window   interfaces.UIControl
	testList *control.UIScrollView
)

// NewScene ...
func NewScene(manager *game.Manager) *Scene {
	s := &Scene{
		children: []interfaces.UIControl{},
		manager:  manager,
	}

	return s
}

// Label ...
func (o *Scene) Label() string {
	return "menu"
}

// TransitionTo ...
func (o *Scene) TransitionTo(t enum.SceneEnum) {
	o.manager.TransitionTo(t)
}

// GetControls ...
func (o *Scene) GetControls() []interfaces.UIControl {
	ret := []interfaces.UIControl{}
	for _, child := range o.children {
		ret = append(ret, child.GetControls()...)
	}
	return ret
}

// DidLoad ...
func (o *Scene) DidLoad() {
	log.Printf("map.DidLoad")
	f := NewWorldMap(o)
	o.children = append(o.children, f)

	window = NewInfoLayer(o)
	f.AppendChild(window)

	testList = NewScrollView(o).(*control.UIScrollView)
	window.AppendChild(testList)
}

// DidActive ...
func (o *Scene) DidActive() {
	log.Printf("map.DidActive")

	cols := []interface{}{
		"ID", "カラム名", "3番目",
	}
	data := [][]interface{}{
		{1, "１行目", "あいうえお"},
		{2, "２行目", "かきくけこ"},
		{3, "３行目", "さしすせそ"},
		{4, "４行目", "たちつてと"},
		{5, "５行目", "なにぬねの"},
		{6, "６行目", "はひふへほ"},
		{7, "７行目", "まみむめも"},
		{8, "８行目", "やいゆえよ"},
		{9, "９行目", "らりるれろ"},
		{10, "１０行目", "わをん"},
		{11, "１１行目", "あいうえお"},
		{12, "１２行目", "かきくけこ"},
		{13, "１３行目", "さしすせそ"},
		{14, "１４行目", "たちつてと"},
		{15, "１５行目", "なにぬねの"},
		{16, "１６行目", "はひふへほ"},
		{17, "１７行目", "まみむめも"},
		{18, "１８行目", "やいゆえよ"},
		{19, "１９行目", "らりるれろ"},
		{20, "２０行目", "わをん"},
		{1, "１行目", "あいうえお"},
		{2, "２行目", "かきくけこ"},
		{3, "３行目", "さしすせそ"},
		{4, "４行目", "たちつてと"},
		{5, "５行目", "なにぬねの"},
		{6, "６行目", "はひふへほ"},
		{7, "７行目", "まみむめも"},
		{8, "８行目", "やいゆえよ"},
		{9, "９行目", "らりるれろ"},
		{10, "１０行目", "わをん"},
		{11, "１１行目", "あいうえお"},
		{12, "１２行目", "かきくけこ"},
		{13, "１３行目", "さしすせそ"},
		{14, "１４行目", "たちつてと"},
		{15, "１５行目", "なにぬねの"},
		{16, "１６行目", "はひふへほ"},
		{17, "１７行目", "まみむめも"},
		{18, "１８行目", "やいゆえよ"},
		{19, "１９行目", "らりるれろ"},
		{20, "２０行目", "わをん"},
		{1, "１行目", "あいうえお"},
		{2, "２行目", "かきくけこ"},
		{3, "３行目", "さしすせそ"},
		{4, "４行目", "たちつてと"},
		{5, "５行目", "なにぬねの"},
		{6, "６行目", "はひふへほ"},
		{7, "７行目", "まみむめも"},
		{8, "８行目", "やいゆえよ"},
		{9, "９行目", "らりるれろ"},
		{10, "１０行目", "わをん"},
		{11, "１１行目", "あいうえお"},
		{12, "１２行目", "かきくけこ"},
		{13, "１３行目", "さしすせそ"},
		{14, "１４行目", "たちつてと"},
		{15, "１５行目", "なにぬねの"},
		{16, "１６行目", "はひふへほ"},
		{17, "１７行目", "まみむめも"},
		{18, "１８行目", "やいゆえよ"},
		{19, "１９行目", "らりるれろ"},
		{20, "２０行目", "わをん"},
		{1, "１行目", "あいうえお"},
		{2, "２行目", "かきくけこ"},
		{3, "３行目", "さしすせそ"},
		{4, "４行目", "たちつてと"},
		{5, "５行目", "なにぬねの"},
		{6, "６行目", "はひふへほ"},
		{7, "７行目", "まみむめも"},
		{8, "８行目", "やいゆえよ"},
		{9, "９行目", "らりるれろ"},
		{10, "１０行目", "わをん"},
		{11, "１１行目", "あいうえお"},
		{12, "１２行目", "かきくけこ"},
		{13, "１３行目", "さしすせそ"},
		{14, "１４行目", "たちつてと"},
		{15, "１５行目", "なにぬねの"},
		{16, "１６行目", "はひふへほ"},
		{17, "１７行目", "まみむめも"},
		{18, "１８行目", "やいゆえよ"},
		{19, "１９行目", "らりるれろ"},
		{20, "２０行目", "わをん"},
	}

	// ヘッダを設定
	testList.SetHeaderRow(cols)

	// 行を設定
	testList.AppendRows(data, func(row *control.ScrollViewRow) {
		log.Printf("RowClickFunc: %#v", row)
	})

}
