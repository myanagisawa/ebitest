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
	// EventTypeWheel ホイール
	EventTypeWheel
	// EventTypeScroll スクロール
	EventTypeScroll
)

// Name EventTypeEnumの区分値表現を返します
func (e EventTypeEnum) Name() string {
	switch e {
	case EventTypeNone:
		return "イベント未検知"
	case EventTypeClick:
		return "クリック"
	case EventTypeFocus:
		return "フォーカス"
	case EventTypeBlur:
		return "ブラー"
	case EventTypeDragging:
		return "ドラッグ中"
	case EventTypeDragDrop:
		return "ドラッグ&ドロップ"
	case EventTypeLongPress:
		return "長押し"
	case EventTypeLongPressReleased:
		return "長押し完了"
	case EventTypeWheel:
		return "ホイール"
	case EventTypeScroll:
		return "スクロール"
	default:
		return "不明"
	}
}
