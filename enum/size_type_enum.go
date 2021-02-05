package enum

type (
	// SizeTypeEnum age
	SizeTypeEnum int
)

const (
	// TypeOriginal オリジナル値
	TypeOriginal SizeTypeEnum = iota + 1
	// TypeScaled スケール適用値
	TypeScaled
)
