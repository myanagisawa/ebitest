package enum

type (
	// ValueTypeEnum age
	ValueTypeEnum int
)

const (
	// TypeLocal ローカル値
	TypeLocal ValueTypeEnum = iota + 1
	// TypeGlobal グローバル値
	TypeGlobal
)
