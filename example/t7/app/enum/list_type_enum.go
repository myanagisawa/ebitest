package enum

type (
	// ListTypeEnum リスト種別
	ListTypeEnum int
)

const (
	// ListTypeUpdate Update処理
	ListTypeUpdate ListTypeEnum = iota
	// ListTypeDraw Draw処理
	ListTypeDraw
	// ListTypeCursor カーソル対象
	ListTypeCursor
)
