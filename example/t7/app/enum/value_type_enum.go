package enum

type (
	// ValueTypeEnum age
	ValueTypeEnum int
)

const (
	// TypeRaw そのまま
	TypeRaw ValueTypeEnum = iota
	// TypeLocal ローカル値
	TypeLocal
	// TypeGlobal グローバル値
	TypeGlobal
)
