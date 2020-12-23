package enum

type (
	// EventTypeEnum イベント種別
	EventTypeEnum int
)

const (
	// EventTypeClick クリック
	EventTypeClick EventTypeEnum = iota
	// EventTypeFocus フォーカス
	EventTypeFocus
	// EventTypeBlur ブラー
	EventTypeBlur
	// EventTypeDrag ドラッグ
	EventTypeDrag
	// EventTypeLongPress 長押し
	EventTypeLongPress
)
