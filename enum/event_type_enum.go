package enum

type (
	// EventTypeEnum イベント種別
	EventTypeEnum int
)

const (
	// EventTypeNone イベント未検知
	EventTypeNone EventTypeEnum = iota
	// EventTypeClick クリック
	EventTypeClick
	// EventTypeFocus フォーカス
	EventTypeFocus
	// EventTypeBlur ブラー
	EventTypeBlur
	// EventTypeDragging ドラッグ中
	EventTypeDragging
	// EventTypeDragDrop ドラッグ&ドロップ
	EventTypeDragDrop
	// EventTypeLongPress 長押し
	EventTypeLongPress
	// EventTypeLongPressReleased 長押し完了
	EventTypeLongPressReleased
)
