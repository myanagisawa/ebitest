package enum

type (
	// FuncTypeEnum データ区分
	FuncTypeEnum int
)

const (
	// FuncTypeDidLoad ...
	FuncTypeDidLoad FuncTypeEnum = iota + 1
	// FuncTypeDidActive ...
	FuncTypeDidActive
)
